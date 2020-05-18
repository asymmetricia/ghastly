package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"
	"strings"
	"sync"
)

// Client is the object used to interface with a homeassistant server. It should not be copied after being created.
type Client struct {
	Token        string
	Server       string
	connection   *websocket.Conn
	connectionMu sync.Mutex
	messageIndex int
}

// RawJSON sends a request using the given method (e.g., GET, POST, DELETE) to the given nominal path. Parameters
// describe URL parameters that are added to the URL, and may be `nil`. body is an object that's converted to JSON and
// supplied as the request body. No request body is supplied if `body` is nil.
//
// Returns the generally JSON-decoded response object, or error, if something happens. (Including, e.g., non-2XX
// responses or a body that isn't parseable JSON).
func (c *Client) RawJSON(method string, path string, parameters map[string]interface{}, body interface{}) (interface{}, error) {
	var rdr io.Reader
	if body != nil {
		j, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("could not marshal body %T as JSON: %v", body, err)
		}
		rdr = bytes.NewBuffer(j)
	}
	return c.Raw(method, path, parameters, rdr)
}

// Raw is as RawJSON, above, except the body is expected to already be an io.Reader.
func (c *Client) Raw(method string, path string, parameters map[string]interface{}, body io.Reader) (interface{}, error) {
	path = strings.Trim(path, "/")
	if !strings.HasPrefix(path, "api/") {
		path = "api/" + path
	}
	haUrl, err := url.Parse(c.Server)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as URL: %v", c.Server, err)
	}
	haUrl.Path += path
	for k, v := range parameters {
		vs, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("REST requests only accept string-value parameters, but %q was %T", k, v)
		}
		haUrl.Query().Set(k, vs)
	}

	req, err := http.NewRequest(method, haUrl.String(), body)
	if err != nil {
		return nil, fmt.Errorf("preparing %q request: %v", method, err)
	}

	return c.executeAndParse(req)
}

// Delete renders `body` as JSON and posts it to the given path.
func (c *Client) Delete(path string, body interface{}) (interface{}, error) {
	return c.RawJSON("DELETE", path, nil, body)
}

// Post renders `body` as JSON and posts it to the given path.
func (c *Client) Post(path string, body interface{}) (interface{}, error) {
	return c.RawJSON("POST", path, nil, body)
}

func (c *Client) Get(path string, parameters map[string]interface{}) (interface{}, error) {
	return c.Raw("GET", path, parameters, nil)
}

func (c *Client) executeAndParse(req *http.Request) (interface{}, error) {
	req.Header.Add("authorization", "Bearer "+c.Token)
	req.Header.Add("content-type", "application/json")

	reqBytes, err := httputil.DumpRequest(req, true)
	if err != nil {
		panic("could not dump request: " + err.Error())
	}
	for _, line := range strings.Split(string(reqBytes), "\n") {
		logrus.Tracef("sent: %q", strings.TrimSpace(line))
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending GET request: %v", err)
	}
	defer res.Body.Close()

	resBytes, err := httputil.DumpResponse(res, true)
	if err != nil {
		panic("could not dump response: " + err.Error())
	}
	for _, line := range strings.Split(string(resBytes), "\n") {
		logrus.Tracef("recv: %q", strings.TrimSpace(line))
	}

	if res.StatusCode/100 != 2 {
		bb, _ := ioutil.ReadAll(res.Body)
		return &ResultMessage{Success: false, Error: ResultError{
			Code:    res.Status,
			Message: string(bb),
		}}, nil
	}

	dec := json.NewDecoder(res.Body)

	var ret []interface{}
	for {
		var obj interface{}
		err := dec.Decode(&obj)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading from server: %v", err)
		}
		ret = append(ret, obj)
	}

	if len(ret) == 0 {
		return nil, nil
	}
	if len(ret) == 1 {
		return &ResultMessage{
			Success: true,
			Result:  ret[0],
		}, nil
	}
	return &ResultMessage{
		Success: true,
		Result:  ret,
	}, nil
}

func (c *Client) connect() error {
	if c.connection != nil {
		return nil
	}

	haUrl, err := url.Parse(c.Server)
	if err != nil {
		return fmt.Errorf("parsing %q: %v", c.Server, err)
	}

	if haUrl.Scheme == "http" {
		haUrl.Scheme = "ws"
	} else {
		haUrl.Scheme = "wss"
	}

	haUrl.Path = "api/websocket"

	c.connection, _, err = websocket.DefaultDialer.Dial(haUrl.String(), nil)
	if err != nil {
		return fmt.Errorf("opening websocket to %q: %v", haUrl.Host, err)
	}

	logrus.Debug("expecting auth_required...")
	msg, err := c.receive()
	if err != nil {
		return fmt.Errorf("handshaking: %v", err)
	}
	if _, ok := msg.(*AuthRequiredMessage); !ok {
		return fmt.Errorf("message was %T, not *AuthRequiredMessage", msg)
	}

	err = c.send(AuthMessage{AccessToken: c.Token})
	if err != nil {
		return fmt.Errorf("authenticating: %v", err)
	}

	msg, err = c.receive()
	switch m := msg.(type) {
	case *AuthOkMessage:
		return nil
	case *AuthInvalidMessage:
		return fmt.Errorf("authentication failed: %v", m.Message)
	default:
		return fmt.Errorf("unexpected response type %T", msg)
	}
}

// receive receives a single message via the websocket and returns the parsed result, or an error if it could not be
// received or not parsed.
func (c *Client) receive() (Message, error) {
	_, data, err := c.connection.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("reading: %v", err)
	}
	logrus.Tracef("    got: %s", string(data))
	ret, err := MessageFromJSON(data)
	if err != nil {
		return nil, fmt.Errorf("parsing: %v", err)
	}
	logrus.Debugf("    got: %s", ret.Type())
	return ret, nil
}

func (c *Client) send(send Message) error {
	if send == nil {
		return fmt.Errorf("cannot send nil message")
	}
	logrus.Debugf("sending: %s", send.Type())
	data, err := json.Marshal(send)
	if err != nil {
		return fmt.Errorf("marshaling frame: %w", err)
	}

	// don't @ me
	var obj map[string]interface{}
	_ = json.Unmarshal(data, &obj)
	obj["type"] = send.Type()
	if _, ok := send.(AuthMessage); !ok {
		c.messageIndex++
		obj["id"] = c.messageIndex
	}
	data, _ = json.Marshal(obj)

	logrus.Tracef("sending: %s", string(data))
	if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
		return fmt.Errorf("sending: %w", err)
	}
	return nil
}

// Exchange exchanges the 'send' Message for a response. It is safe for use by multiple goroutines.
func (c *Client) Exchange(send Message) (Message, error) {
	c.connectionMu.Lock()
	defer c.connectionMu.Unlock()

	// establish a new connection if needed
	if c.connection == nil {
		if err := c.connect(); err != nil {
			return nil, fmt.Errorf("connecting: %v", err)
		}
	}

	if err := c.send(send); err != nil {
		return nil, err
	}

	return c.receive()
}

type Message interface {
	Type() string
}

func RegisterMessageType(typ ...Message) {
	for _, t := range typ {
		MessageHandlers[t.Type()] = t
	}
}

var MessageHandlers = map[string]Message{}

// MessageFromJSON returns one of the various ___Message structs described
// above, or an error if we can't do so.
func MessageFromJSON(data []byte) (Message, error) {
	var dataMap map[string]interface{}
	if err := json.Unmarshal(data, &dataMap); err != nil {
		return nil, err
	}
	typ, typOk := dataMap["type"].(string)
	if !typOk {
		return nil, errors.New("object did not contain string Type field")
	}
	proto, ok := MessageHandlers[typ]
	if !ok {
		return nil, fmt.Errorf("unhandled message type %q", typ)
	}
	ret := reflect.New(reflect.TypeOf(proto)).Interface()
	if err := json.Unmarshal(data, &ret); err != nil {
		return nil, fmt.Errorf("parsing %q into %T: %v", typ, ret, err)
	}
	return ret.(Message), nil
}
