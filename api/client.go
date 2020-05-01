package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Client struct {
	Token      string
	Server     string
	connection *websocket.Conn
}

func (c *Client) Post(path string, body io.Reader) (interface{}, error) {
	path = strings.Trim(path, "/")
	if !strings.HasPrefix(path, "api/") {
		path = "api/" + path
	}
	haUrl, err := url.Parse(c.Server)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as URL: %v", c.Server, err)
	}
	haUrl.Path += path

	req, err := http.NewRequest("POST", haUrl.String(), body)
	if err != nil {
		return nil, fmt.Errorf("preparing POST request: %v", err)
	}

	return c.executeAndParse(req)
}

func (c *Client) Get(path string, parameters map[string]string) (interface{}, error) {
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
		haUrl.Query().Set(k, v)
	}

	req, err := http.NewRequest("GET", haUrl.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("preparing GET request: %v", err)
	}

	return c.executeAndParse(req)
}

func (c *Client) executeAndParse(req *http.Request) (interface{}, error) {
	req.Header.Add("authorization", "Bearer "+c.Token)
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending GET request: %v", err)
	}
	defer res.Body.Close()

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
	msg, err := c.Exchange(nil)
	if err != nil {
		return fmt.Errorf("handshaking: %v", err)
	}
	if _, ok := msg.(*AuthRequiredMessage); !ok {
		return fmt.Errorf("message was %T, not *AuthRequiredMessage", msg)
	}

	msg, err = c.Exchange(AuthMessage{AccessToken: c.Token})
	if err != nil {
		return fmt.Errorf("authenticating: %v", err)
	}
	switch m := msg.(type) {
	case *AuthOkMessage:
		return nil
	case *AuthInvalidMessage:
		return fmt.Errorf("authentication failed: %v", m.Message)
	default:
		return fmt.Errorf("unexpected response type %T", msg)
	}
}

func (c *Client) Exchange(send Message) (Message, error) {
	if c.connection == nil {
		if err := c.connect(); err != nil {
			return nil, fmt.Errorf("connecting: %v", err)
		}
	}

	if send != nil {
		logrus.Debugf("sending: %s", send.Type())
		data, err := json.Marshal(send)
		if err != nil {
			return nil, fmt.Errorf("marshaling frame: %v", err)
		}
		// don't @ me
		var obj map[string]interface{}
		_ = json.Unmarshal(data, &obj)
		obj["type"] = send.Type()
		if _, ok := send.(AuthMessage); !ok {
			obj["id"] = 1
		}
		data, _ = json.Marshal(obj)

		logrus.Tracef("sending: %s", string(data))
		if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
			return nil, fmt.Errorf("sending: %v", err)
		}
	}

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
