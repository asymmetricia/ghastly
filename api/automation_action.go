package api

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type AutomationAction struct {
	Action
}

func (a *AutomationAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.Action)
}

var _ json.Marshaler = (*AutomationAction)(nil)

type Action interface {
	KeyFieldName() string
}

type EventAction struct {
	Event             string                 `json:"event"`
	EventData         map[string]interface{} `json:"event_data,omitempty"`
	EventDataTemplate map[string]interface{} `json:"event_data_template,omitempty"`
}

func (*EventAction) KeyFieldName() string { return "event" }

type ServiceAction struct {
	Service  string                 `json:"service"`
	EntityId *string                `json:"entity_id,omitempty"`
	Data     map[string]interface{} `json:"data,omitempty"`
}

func (*ServiceAction) KeyFieldName() string { return "service" }

type DeviceAction struct {
	DeviceId string `json:"device_id"`
	Domain   string `json:"domain"`
	EntityId string `json:"entity_id"`
}

func (*DeviceAction) KeyFieldName() string { return "device_id" }

type DelayAction struct {
	Delay int `json:"delay"`
}

func (*DelayAction) KeyFieldName() string { return "delay" }

type SceneAction struct {
	Scene string `json:"scene"`
}

func (*SceneAction) KeyFieldName() string { return "scene" }

type WaitAction struct {
	WaitTemplate string `json:"wait_template"`
	Timeout      int    `json:"timeout,omitempty"`
}

func (*WaitAction) KeyFieldName() string { return "wait_template" }

func (a *AutomationAction) UnmarshalJSON(data []byte) error {
	var generic map[string]interface{}
	if err := json.Unmarshal(data, &generic); err != nil {
		return err
	}

	if len(generic) == 0 {
		return nil
	}

	var prototype Action
	for _, candidate := range []Action{
		(*EventAction)(nil),
		(*ServiceAction)(nil),
		(*DeviceAction)(nil),
		(*DelayAction)(nil),
		(*SceneAction)(nil),
		(*WaitAction)(nil),
	} {
		if _, ok := generic[candidate.KeyFieldName()]; ok {
			prototype = candidate
			break
		}
	}

	if prototype == nil {
		return fmt.Errorf("action %q was not a recognized type", string(data))
	}

	typ := reflect.TypeOf(prototype).Elem()
	action := reflect.New(typ)

	for fieldIdx := 0; fieldIdx < typ.NumField(); fieldIdx++ {
		targetField := action.Elem().Field(fieldIdx)
		fieldName := strings.Split(typ.Field(fieldIdx).Tag.Get("json"), ",")[0]
		valueToAssign, ok := generic[fieldName]
		if !ok {
			continue
		}

		if reflect.TypeOf(valueToAssign) != targetField.Type() {
			return fmt.Errorf("%q expects %s for field %s, but input had %s", prototype.KeyFieldName(),
				reflect.TypeOf(targetField), fieldName, reflect.TypeOf(valueToAssign))
		}

		targetField.Set(reflect.ValueOf(valueToAssign))
	}

	a.Action = action.Interface().(Action)
	return nil
}

var _ json.Unmarshaler = (*AutomationAction)(nil)
