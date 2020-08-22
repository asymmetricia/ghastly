package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"

	"github.com/sirupsen/logrus"
)

type domain struct {
	Domain   string
	Services map[string]Service
}

type Service struct {
	client      *Client
	Domain      string                   `json:"domain"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Fields      map[string]*ServiceField `json:"fields"`
}

type ServiceField struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Type        ServiceFieldType `json:"type"`
	Example     interface{}      `json:"example"`
	Values      []interface{}    `json:"values"`
	Default     interface{}      `json:"default"`
}

type ServiceFieldType string

const (
	Number  ServiceFieldType = "number"
	String                   = "string"
	Boolean                  = "boolean"
	Values                   = "values"
)

func (s *ServiceField) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewBuffer(data))
	dec.UseNumber()

	const (
		stateInit = iota
		stateKey
		stateDescription
		stateExample
		stateDefault
		stateEnd
	)

	state := 0
loop:
	for {
		t, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("decoding service field %q: %w", string(data), err)
		}
		switch state {
		case stateInit:
			delim, ok := t.(json.Delim)
			if !ok || delim != '{' {
				return fmt.Errorf("decoding service field %q: expected { but "+
					"got %v", string(data), t)
			}
			state = stateKey
		case stateKey:
			key, ok := t.(string)
			if !ok {
				return fmt.Errorf("decoding service field %q: expected "+
					"string but got %v", string(data), t)
			}
			if key == "description" {
				state = stateDescription
			} else if key == "example" || key == "exampl" /* for fuck's sake, homeassistant */ {
				state = stateExample
			} else if key == "values" {
				s.Type = Values
				if err := dec.Decode(&(s.Values)); err != nil {
					return fmt.Errorf("decoding service field %q: expected list "+
						"of strings but: %w", string(data), err)
				}
			} else if key == "default" {
				state = stateDefault
			} else {
				return fmt.Errorf("unexpected key value %q in service field "+
					"%q", key, string(data))
			}
		case stateDescription:
			desc, ok := t.(string)
			if !ok {
				return fmt.Errorf("decoding service field %q: expected "+
					"string for description value, but got %T %v", string(data), t, t)
			}
			s.Description = desc
			state = stateKey
		case stateExample:
			switch v := t.(type) {
			case string:
				s.Type = String
				s.Example = v
			case json.Number:
				s.Type = Number
				s.Example, _ = v.Float64()
			case bool:
				s.Type = Boolean
				s.Example = v
			default:
				return fmt.Errorf("decoding service field %q: unepected example type %T", string(data), t)
			}
			state = stateKey
		case stateDefault:
			switch v := t.(type) {
			case string:
				s.Type = String
				s.Default = v
			case json.Number:
				s.Type = Number
				s.Default, _ = v.Float64()
			case bool:
				s.Type = Boolean
				s.Default = v
			default:
				return fmt.Errorf("decoding service field %q: unepected example type %T", string(data), t)
			}
			state = stateKey
		case stateEnd:
			delim, ok := t.(json.Delim)
			if !ok || delim != '}' {
				return fmt.Errorf("decoding service field %q: expected } but got %c", string(data), delim)
			}
			break loop
		}
		if !dec.More() {
			state = stateEnd
		}
	}

	return nil
}

//var _ json.Marshaler = (*ServiceField)(nil)
var _ json.Unmarshaler = (*ServiceField)(nil)

func (s ServiceFieldType) String() string {
	return string(s)
}

func (c *Client) ListServices() ([]Service, error) {
	servicesI, err := c.RawRESTGetAs("services", nil, []domain{})
	if err != nil {
		return nil, err
	}

	var ret []Service
	for _, domain := range servicesI.([]domain) {
		for name, svc := range domain.Services {
			svc.client = c
			svc.Domain = domain.Domain
			svc.Name = name
			for name, f := range svc.Fields {
				f.Name = name
			}
			ret = append(ret, svc)
		}
	}

	return ret, nil
}

var ServiceNotFound = errors.New("not found")

// GetService returns the service with the given domain and service. If the
// service can't be retrieved, the returned service will be nil and the error
// will be non-nil. If the service does not exist but no other error occurs,
// an error wrapping ServiceNotFound will be returned.
func (c *Client) GetService(domain, service string) (*Service, error) {
	svcs, err := c.ListServices()
	if err != nil {
		return nil, fmt.Errorf("retrieving services list: %w", err)
	}

	for _, svc := range svcs {
		if svc.Domain == domain && svc.Name == service {
			return &svc, nil
		}
	}

	return nil, fmt.Errorf("service %q in domain %q %w", service, domain, ServiceNotFound)
}

func (s *Service) Call(data map[string]interface{}) ([]State, error) {
	for f, v := range data {
		logrus.Tracef("validating %s: %v", f, v)
		field, ok := s.Fields[f]
		if !ok {
			return nil, fmt.Errorf("service %s.%s does not have field %q", s.Domain, s.Name, f)
		}

		switch field.Type {
		case String:
			_, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("service %s.%s expects field %s to be %s, "+
					"but %T was provided", s.Domain, s.Name, f, field.Type, v)
			}
		case Number:
			_, fOk := v.(float64)
			_, iOk := v.(int)
			if !fOk && !iOk {
				return nil, fmt.Errorf("service %s.%s expects field %s to be int "+
					"or float64, but %T was provided", s.Domain, s.Name, f, v)
			}
		case Values:
			found := false
			for _, value := range field.Values {
				if reflect.TypeOf(value) == reflect.TypeOf(v) && reflect.ValueOf(value) == reflect.ValueOf(v) {
					found = true
					break
				}
			}
			if !found {
				vs, _ := json.Marshal(field.Values)
				return nil, fmt.Errorf("service %s.%s field %s only allows "+
					"values %s; %q was not found", s.Domain, s.Name, f, vs, v)
			}
		case Boolean:
			if _, ok := v.(bool); !ok {
				return nil, fmt.Errorf("service %s.%s expects field %s to be %s, "+
					"but %T was provided", s.Domain, s.Name, f, field.Type, v)
			}
		}
	}

	res, err := s.client.Post(fmt.Sprintf("/api/services/%s/%s", s.Domain,
		s.Name), data)

	if err != nil {
		return nil, fmt.Errorf("during POST: %w", err)
	}

	logrus.Debugf("%+v", res)
	resObj, ok := res.(*ResultMessage)

	if !ok {
		return nil, fmt.Errorf("expected ResultMessage back from Post but got %T", res)
	}

	if !resObj.Success {
		return nil, fmt.Errorf("service call failed: %+v", resObj.Error)
	}

	resultJson, err := json.Marshal(resObj.Result)
	if err != nil {
		return nil, fmt.Errorf("could not marshal result %+v as part of "+
			"converting to []State: %w", resObj.Result, err)
	}
	var ret []State
	if err := json.Unmarshal(resultJson, &ret); err != nil {
		return nil, fmt.Errorf("could not unmarshal result JSON %s into "+
			"[]State: %w", resultJson, err)
	}

	return ret, nil
}
