package api

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type AutomationCondition struct{ Condition }

func (a *AutomationCondition) UnmarshalJSON(data []byte) error {
	var generic map[string]interface{}
	if err := json.Unmarshal(data, &generic); err != nil {
		return err
	}
	condition, err := conditionFromMap(generic)
	if err != nil {
		return err
	}

	a.Condition = condition
	return nil
}

func (a AutomationCondition) MarshalJSON() ([]byte, error) {
	return json.Marshal(conditionToMap(a.Condition))
}

var _ json.Unmarshaler = (*AutomationCondition)(nil)
var _ json.Marshaler = (*AutomationCondition)(nil)

func conditionToMap(c Condition) map[string]interface{} {
	ret := map[string]interface{}{}
	ret["condition"] = c.ConditionKey()
	typ := reflect.TypeOf(c).Elem()
	val := reflect.ValueOf(c).Elem()
	for fieldIdx := 0; fieldIdx < typ.NumField(); fieldIdx++ {
		name := strings.Split(typ.Field(fieldIdx).Tag.Get("json"), ",")[0]
		switch v := val.Field(fieldIdx).Interface().(type) {
		case []Condition:
			var list []map[string]interface{}
			for _, c := range v {
				list = append(list, conditionToMap(c))
			}
			ret[name] = list
		default:
			ret[name] = v
		}
	}

	return ret
}

func conditionFromMap(generic map[string]interface{}) (Condition, error) {
	conditionName, ok := generic["condition"]
	if !ok {
		return nil, fmt.Errorf("condition %v did not have `condition` key", generic)
	}

	var prototype Condition
candidates:
	for _, candidate := range []Condition{
		(*AndCondition)(nil),
		(*OrCondition)(nil),
		(*NotCondition)(nil),
		(*StateCondition)(nil),
		(*NumericStateCondition)(nil),
		(*SunCondition)(nil),
		(*ZoneCondition)(nil),
		(*TimeCondition)(nil),
		(*TemplateCondition)(nil),
	} {
		if conditionName == candidate.ConditionKey() {
			prototype = candidate
			break candidates
		}
	}
	if prototype == nil {
		return nil, fmt.Errorf("no match for condition %q", conditionName)
	}

	typ := reflect.TypeOf(prototype).Elem()
	condition := reflect.New(typ)

	for fieldIdx := 0; fieldIdx < typ.NumField(); fieldIdx++ {
		fieldName := strings.Split(typ.Field(fieldIdx).Tag.Get("json"), ",")[0]
		valueToAssign, ok := generic[fieldName]
		if !ok {
			continue
		}
		targetField := condition.Elem().Field(fieldIdx)
		switch targetField.Interface().(type) {
		case []Condition:
			conditions, ok := valueToAssign.([]interface{})
			if !ok {
				return nil, fmt.Errorf("%s attribute was %T, not []interface{}", fieldName, valueToAssign)
			}

			for conditionIdx, condI := range conditions {
				condGeneric, ok := condI.(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("%s attribute entry %d was %T, not "+
						"map[string]interface{}", fieldName, conditionIdx, condI)
				}

				cond, err := conditionFromMap(condGeneric)
				if err != nil {
					return nil, fmt.Errorf("%s attribute entry %d could not be decoded as "+
						"Condition: %w", fieldName, conditionIdx, err)
				}

				targetField.Set(reflect.Append(targetField, reflect.ValueOf(cond)))
			}
		case StringOrFloat:
			switch v := valueToAssign.(type) {
			case string:
				condition.Elem().Field(fieldIdx).Set(reflect.ValueOf(StringOrIntFromString(v)))
			case float64:
				condition.Elem().Field(fieldIdx).Set(reflect.ValueOf(StringOrIntFromFloat64(v)))
			default:
				return nil, fmt.Errorf("%s attribute was %T, not string or int", fieldName, valueToAssign)
			}
		default:
			if reflect.TypeOf(valueToAssign) != targetField.Type() {
				return nil, fmt.Errorf("%q expects %s for field %s, but input had %s", conditionName,
					reflect.TypeOf(targetField), fieldName, reflect.TypeOf(valueToAssign))
			}
			targetField.Set(reflect.ValueOf(valueToAssign))
		}
	}

	return condition.Interface().(Condition), nil
}

type Condition interface {
	ConditionKey() string
}

type AndCondition struct {
	Conditions []Condition `json:"conditions"`
}

func (*AndCondition) ConditionKey() string { return "and" }

type OrCondition struct {
	Conditions []Condition `json:"conditions"`
}

func (*OrCondition) ConditionKey() string { return "or" }

type NotCondition struct {
	Conditions []Condition `json:"conditions"`
}

func (*NotCondition) ConditionKey() string { return "not" }

type StateCondition struct {
	EntityId string        `json:"entity_id"`
	State    StringOrFloat `json:"state"`
}

func (*StateCondition) ConditionKey() string { return "state" }

type NumericStateCondition struct {
	EntityId string  `json:"entity_id"`
	Above    float64 `json:"above"`
}

func (*NumericStateCondition) ConditionKey() string { return "numeric_state" }

type SunCondition struct {
	AfterOffset  float64 `json:"after_offset"`
	BeforeOffset float64 `json:"before_offset"`
	After        string  `json:"after"`
	Before       string  `json:"before"`
}

func (*SunCondition) ConditionKey() string { return "sun" }

type ZoneCondition struct {
	EntityId string `json:"entity_id"`
	Zone     string `json:"zone"`
}

func (*ZoneCondition) ConditionKey() string { return "zone" }

type TimeCondition struct {
	After  string `json:"after"`
	Before string `json:"before"`
}

func (*TimeCondition) ConditionKey() string { return "time" }

type TemplateCondition struct {
	ValueTemplate string `json:"value_template"`
}

func (*TemplateCondition) ConditionKey() string { return "template" }

// StringOrFloat represents a value that is either a string or an integer.
type StringOrFloat struct {
	*string
	*float64
}

func (s StringOrFloat) String() string {
	if s.string == nil && s.float64 == nil {
		return ""
	}

	if s.string == nil {
		return fmt.Sprintf("%f", *s.float64)
	}

	return *s.string
}

func StringOrIntFromString(s string) StringOrFloat {
	ret := StringOrFloat{string: new(string)}
	*ret.string = s
	return ret
}

func StringOrIntFromFloat64(f float64) StringOrFloat {
	ret := StringOrFloat{float64: new(float64)}
	*ret.float64 = f
	return ret
}

func (s StringOrFloat) UnmarshalJSON(data []byte) error {
	var something interface{}
	if err := json.Unmarshal(data, &something); err != nil {
		return err
	}

	switch v := something.(type) {
	case string:
		s.string = new(string)
		*s.string = v
		return nil
	case float64:
		s.float64 = new(float64)
		*s.float64 = v
		return nil
	}

	return fmt.Errorf("expecting string or int but got %T", something)
}

func (s StringOrFloat) MarshalJSON() ([]byte, error) {
	if s.string == nil && s.float64 == nil {
		return []byte(`""`), nil
	}

	if s.string != nil {
		return json.Marshal(s.string)
	}

	return json.Marshal(s.float64)
}

var _ json.Marshaler = (*StringOrFloat)(nil)
var _ json.Unmarshaler = (*StringOrFloat)(nil)
