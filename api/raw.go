package api

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// RawWebsocketRequest exchanges the given message and returns the Result.Result if everything goes well. If anything
// else happens (failure communicating, rejected request, etc) the result will be nil and the error wll be non-nil.
// Note that the returned interface will be generically typed (e.g., maps & slices). See RawWebsocketRequestAs.
func (c *Client) RawWebsocketRequest(message Message) (interface{}, error) {
	return process(c.Exchange(message))
}

// RawWebsocketRequest exchanges the given message and returns the Result.Result if everything goes well. If anything
// else happens (failure communicating, rejected request, etc) the result will be nil and the error wll be non-nil.
// The returned interface will be the same type as the given prototype. If conversion cannot be achieved, an error will
// be returned.
func (c *Client) RawWebsocketRequestAs(message Message, prototype interface{}) (interface{}, error) {
	obj, err := process(c.Exchange(message))
	return convert(obj, prototype, err)
}

func convert(obj interface{}, prototype interface{}, err error) (interface{}, error) {
	if err != nil {
		return nil, err
	}

	ob, _ := json.Marshal(obj)
	ret := reflect.New(reflect.TypeOf(prototype)).Interface()
	if err := json.Unmarshal(ob, &ret); err != nil {
		return nil, fmt.Errorf("could not convert response to %T: %v", ret, err)
	}

	return reflect.Indirect(reflect.ValueOf(ret)).Interface(), nil
}

func process(resultInterface interface{}, err error) (interface{}, error) {
	if err != nil {
		return nil, fmt.Errorf("excanging message: %v", err)
	}

	result, ok := resultInterface.(*ResultMessage)
	if !ok {
		return nil, fmt.Errorf("response was %T, not ResultMessage", resultInterface)
	}

	if !result.Success {
		return nil, fmt.Errorf("result was not successful: %+v", result.Error)
	}

	return result.Result, nil
}

// RawRESTGet requests the given path via the REST API and returns the JSON-decoded request body if everything goes
// well. If anything else happens (failure communicating, rejected request, etc.) the result will be nil and the error
// will be non-nil.
// Note that the returned interface will be generically typed (e.g., maps & slices). See RawRESTRequestAs.
func (c *Client) RawRESTGet(path string, parameters map[string]interface{}) (interface{}, error) {
	return process(c.Get(path, parameters))
}

// RawRESTGetAs requests the given path via the REST API and returns the JSON-decoded request body if everything goes
// well. If anything else happens (failure communicating, rejected request, etc.) the result will be nil and the error
// will be non-nil.
// The returned interface will be the same type as the given prototype. If conversion cannot be achieved, an error will
// be returned.
func (c *Client) RawRESTGetAs(path string, parameters map[string]interface{}, prototype interface{}) (interface{}, error) {
	obj, err := process(c.Get(path, parameters))
	return convert(obj, prototype, err)
}

// RawRESTPost requests the given path via the REST API and returns the JSON-decoded request body if everything goes
// well. If anything else happens (failure communicating, rejected request, etc.) the result will be nil and the error
// will be non-nil.
// Note that the returned interface will be generically typed (e.g., maps & slices). See RawRESTRequestAs.
func (c *Client) RawRESTPost(path string, parameters map[string]interface{}) (interface{}, error) {
	return process(c.Post(path, parameters))
}

// RawRESTPostAs requests the given path via the REST API and returns the JSON-decoded request body if everything goes
// well. If anything else happens (failure communicating, rejected request, etc.) the result will be nil and the error
// will be non-nil.
// The returned interface will be the same type as the given prototype. If conversion cannot be achieved, an error will
// be returned.
func (c *Client) RawRESTPostAs(path string, parameters map[string]interface{}, prototype interface{}) (interface{}, error) {
	obj, err := process(c.Post(path, parameters))
	return convert(obj, prototype, err)
}

// RawRESTDelete requests the given path via the REST API and returns the JSON-decoded request body if everything goes
// well. If anything else happens (failure communicating, rejected request, etc.) the result will be nil and the error
// will be non-nil.
// Note that the returned interface will be generically typed (e.g., maps & slices). See RawRESTRequestAs.
func (c *Client) RawRESTDelete(path string, parameters map[string]interface{}) (interface{}, error) {
	return process(c.Delete(path, parameters))
}

// RawRESTDeleteAs requests the given path via the REST API and returns the JSON-decoded request body if everything goes
// well. If anything else happens (failure communicating, rejected request, etc.) the result will be nil and the error
// will be non-nil.
// The returned interface will be the same type as the given prototype. If conversion cannot be achieved, an error will
// be returned.
func (c *Client) RawRESTDeleteAs(path string, parameters map[string]interface{}, prototype interface{}) (interface{}, error) {
	obj, err := process(c.Delete(path, parameters))
	return convert(obj, prototype, err)
}
