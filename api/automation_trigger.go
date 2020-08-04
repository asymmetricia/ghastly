package api

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type AutomationTrigger struct {
	Trigger
}

func (a AutomationTrigger) MarshalJSON() ([]byte, error) {
	return json.Marshal(triggerToMap(a.Trigger))
}

func triggerToMap(t Trigger) map[string]interface{} {
	ret := map[string]interface{}{
		"platform": t.Platform(),
	}
	typ := reflect.TypeOf(t).Elem()
	val := reflect.ValueOf(t).Elem()
	for fieldIdx := 0; fieldIdx < typ.NumField(); fieldIdx++ {
		name := strings.Split(typ.Field(fieldIdx).Tag.Get("json"), ",")[0]
		ret[name] = val.Field(fieldIdx).Interface()
	}

	return ret
}

func (a *AutomationTrigger) UnmarshalJSON(data []byte) error {
	var generic map[string]interface{}
	if err := json.Unmarshal(data, &generic); err != nil {
		return err
	}
	trigger, err := triggerFromMap(generic)
	if err != nil {
		return err
	}

	a.Trigger = trigger
	return nil
}

var _ json.Unmarshaler = (*AutomationTrigger)(nil)
var _ json.Marshaler = AutomationTrigger{}

type Trigger interface {
	Platform() string
}

func triggerFromMap(generic map[string]interface{}) (Trigger, error) {
	platformName, ok := generic["platform"]
	if !ok {
		return nil, fmt.Errorf("trigger %v did not have `platform` key", generic)
	}

	var prototype Trigger
candidates:
	for _, candidate := range []Trigger{
		(*StateTrigger)(nil),
		(*MqttTrigger)(nil),
		(*GeoLocationTrigger)(nil),
		(*HassTrigger)(nil),
		(*NumericStateTrigger)(nil),
		(*SunTrigger)(nil),
		(*TimePatternTrigger)(nil),
		(*WebhookTrigger)(nil),
		(*ZoneTrigger)(nil),
		(*TimeTrigger)(nil),
		(*TemplateTrigger)(nil),
		(*EventTrigger)(nil),
		(*DeviceTrigger)(nil),
	} {
		if platformName == candidate.Platform() {
			prototype = candidate
			break candidates
		}
	}
	if prototype == nil {
		return nil, fmt.Errorf("no match for trigger platform %q", platformName)
	}

	typ := reflect.TypeOf(prototype).Elem()
	trigger := reflect.New(typ)

	for fieldIdx := 0; fieldIdx < typ.NumField(); fieldIdx++ {
		targetField := trigger.Elem().Field(fieldIdx)
		fieldName := strings.Split(typ.Field(fieldIdx).Tag.Get("json"), ",")[0]
		valueToAssign, ok := generic[fieldName]
		if !ok {
			continue
		}
		switch trigger.Elem().Field(fieldIdx).Interface().(type) {
		case StringOrFloat:
			switch v := valueToAssign.(type) {
			case string:
				targetField.Set(reflect.ValueOf(StringOrIntFromString(v)))
			case float64:
				targetField.Set(reflect.ValueOf(StringOrIntFromFloat64(v)))
			default:
				return nil, fmt.Errorf("%s attribute was %T, not string or int", fieldName, valueToAssign)
			}
		case time.Duration:
			switch v := valueToAssign.(type) {
			case string:
				dur, err := time.ParseDuration(v)
				if err != nil {
					return nil, fmt.Errorf("%s attribute %q could not be parsed as duration: %w", fieldName, v, err)
				}
				targetField.Set(reflect.ValueOf(dur))
			case float64:
				targetField.Set(reflect.ValueOf(time.Duration(v) * time.Second))
			case time.Duration:
				targetField.Set(reflect.ValueOf(v))
			default:
				return nil, fmt.Errorf("%s attribute was expected to be string, float, or time.Duration, but "+
					"was %T instead", fieldName, v)
			}
		case time.Time:
			switch v := valueToAssign.(type) {
			case time.Time:
				targetField.Set(reflect.ValueOf(v))
			case string:
				t, err := time.Parse(time.RFC3339, v)
				if err != nil {
					return nil, fmt.Errorf("%s attribute %q could not be parsed as time: %w", fieldName, v, err)
				}
				targetField.Set(reflect.ValueOf(t))
			default:
				return nil, fmt.Errorf("%s attribute was expected to be RFC3339 string or time.Time, but "+
					"was %T", fieldName, v)
			}
		case interface{}:
			targetField.Set(reflect.ValueOf(valueToAssign))
		default:
			if reflect.TypeOf(valueToAssign) != targetField.Type() {
				return nil, fmt.Errorf("%q expects %s for field %s, but input had %s", platformName,
					reflect.TypeOf(targetField), fieldName, reflect.TypeOf(valueToAssign))
			}
			targetField.Set(reflect.ValueOf(valueToAssign))
		}
	}

	return trigger.Interface().(Trigger), nil
}

type StateTrigger struct {
	EntityId string        `json:"entity_id,omitempty"`
	From     StringOrFloat `json:"from"`
	To       StringOrFloat `json:"to"`
	For      time.Duration `json:"for"`
}

func (*StateTrigger) Platform() string { return "state" }

type MqttTrigger struct {
	Topic   string `json:"topic"`
	Payload string `json:"payload,omitempty"`
}

func (*MqttTrigger) Platform() string { return "mqtt" }

type GeoLocationTrigger struct {
	Source string `json:"source"`
	Zone   string `json:"zone"`
	Event  string `json:"event"`
}

func (*GeoLocationTrigger) Platform() string { return "geo_location" }

type HassTrigger struct {
	Event string `json:"event"`
}

func (*HassTrigger) Platform() string { return "homeassistant" }

type NumericStateTrigger struct {
	EntityId      string        `json:"entity_id"`
	Above         float64       `json:"above"`
	Below         float64       `json:"below"`
	ValueTemplate string        `json:"value_template"`
	For           time.Duration `json:"for"`
}

func (*NumericStateTrigger) Platform() string { return "numeric_state" }

type SunTrigger struct {
	Offset float64 `json:"offset"`
	Event  string  `json:"event"`
}

func (*SunTrigger) Platform() string { return "sun" }

type TimePatternTrigger struct {
	Hours   StringOrFloat `json:"hours"`
	Minutes StringOrFloat `json:"minutes"`
	Seconds StringOrFloat `json:"seconds"`
}

func (*TimePatternTrigger) Platform() string { return "time_pattern" }

type WebhookTrigger struct {
	WebhookId string `json:"webhook_id"`
}

func (*WebhookTrigger) Platform() string { return "webhook" }

type ZoneTrigger struct {
	EntityId string `json:"entity_id"`
	Zone     string `json:"zone"`
	Event    string `json:"event"`
}

func (*ZoneTrigger) Platform() string { return "zone" }

type TimeTrigger struct {
	// A time like HH:MM:SS, 24-hour time.
	At string `json:"at"`
}

func (*TimeTrigger) Platform() string { return "time" }

type TemplateTrigger struct {
	ValueTemplate string `json:"value_template"`
}

func (*TemplateTrigger) Platform() string { return "template" }

type EventTrigger struct {
	EventType string      `json:"event_type"`
	EventData interface{} `json:"event_data"`
}

func (*EventTrigger) Platform() string { return "event" }

type DeviceTrigger struct {
	DeviceId string `json:"device_id"`
	Domain   string `json:"domain"`
	EntityId string `json:"entity_id"`
	Type     string `json:"type,omitempty"`
	Subtype  string `json:"subtype,omitempty"`
	Event    string `json:"event,omitempty"`
}

func (*DeviceTrigger) Platform() string { return "device" }
