package search

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// ToAttributes uses JSON struct tags to convert the object to a flat map of
// attributes. Nested structures receive names concatenated with `_`.
// Collisions are resolved arbitrarily.
func Attributes(obj interface{}) (map[string]string, error) {
	val := reflect.ValueOf(obj)
	typ := val.Type()
	if typ.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil, nil
		}
		val = reflect.Indirect(val)
		typ = val.Type()
	}

	switch typ.Kind() {
	case reflect.Struct:
		ret := map[string]string{}
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			jsonTag := field.Tag.Get("json")
			if jsonTag == "-" {
				continue
			}
			name := strings.ToLower(field.Name)
			jsonName := strings.Split(jsonTag, ",")[0]
			if jsonName != "" {
				name = jsonName
			}
			attributes, err := Attributes(val.Field(i).Interface())
			if err != nil {
				return nil, fmt.Errorf("%s: %v", field.Name, err)
			}
			for k, v := range attributes {
				if k == "" {
					ret[name] = v
				} else {
					ret[name+"_"+k] = v
				}
			}
		}
		return ret, nil
	case reflect.String:
		return map[string]string{"": val.String()}, nil
	case reflect.Slice:
		fallthrough
	case reflect.Array:
		ret := map[string]string{}
		for i := 0; i < val.Len(); i++ {
			attributes, err := Attributes(val.Index(i).Interface())
			if err != nil {
				return nil, fmt.Errorf("%d: %v", i, err)
			}
			for k, v := range attributes {
				if k == "" {
					ret[strconv.Itoa(i)] = v
				} else {
					ret[strconv.Itoa(i)+"_"+k] = v
				}
			}
		}
		return ret, nil
	default:
		s, ok := obj.(fmt.Stringer)
		if !ok {
			return nil, fmt.Errorf("cannot convert %T to attributes", obj)
		}
		return map[string]string{"": s.String()}, nil
	}
}
