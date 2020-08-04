package api

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConditionFromMap(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Condition
		wantErr bool
	}{
		{
			"state condition decode",
			`{"condition": "state", "entity_id": "test-eid", "state": "test-state"}`,
			&StateCondition{
				EntityId: "test-eid",
				State:    StringOrIntFromString("test-state"),
			},
			false,
		}, {
			"state condition decode numeric",
			`{"condition": "state", "entity_id": "test-eid", "state": 1}`,
			&StateCondition{
				EntityId: "test-eid",
				State:    StringOrIntFromFloat64(1),
			},
			false,
		}, {
			"mismatched type",
			`{"condition": "state", "entity_id": 99999, "state": "test-state"}`,
			nil,
			true,
		},
		{
			"simple logical",
			`{"condition": "and", "conditions": [ {"condition": "state", "entity_id": "test-eid", "state": "test-state"} ] }`,
			&AndCondition{
				Conditions: []Condition{&StateCondition{
					EntityId: "test-eid",
					State:    StringOrIntFromString("test-state"),
				}},
			},
			false,
		},
		{"nested logical",
			`{"condition": "or", "conditions": [` +
				`{"condition": "and", "conditions": [` +
				`{"condition": "state", "entity_id": "test-a", "state": "hi"},` +
				`{"condition": "state", "entity_id": "test-b", "state": 99}` +
				`]}` +
				`]}`,
			&OrCondition{
				Conditions: []Condition{&AndCondition{
					Conditions: []Condition{
						&StateCondition{"test-a", StringOrIntFromString("hi")},
						&StateCondition{"test-b", StringOrIntFromFloat64(99)},
					},
				}},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got *AutomationCondition
			err := json.Unmarshal([]byte(tt.input), &got)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got.Condition)

				// normalize input JSON
				var obj interface{}
				require.NoError(t, json.Unmarshal([]byte(tt.input), &obj))
				expectedJson, err := json.Marshal(obj)
				require.NoError(t, err)

				// marshal *AutomationCondition
				gotJson, err := json.Marshal(got)
				require.NoError(t, err)

				require.Equal(t, string(expectedJson), string(gotJson))
			}
		})
	}
}
