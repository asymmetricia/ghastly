package api

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestServiceField_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{"bool", `{ "description": "New value of axillary heater.", "example": true }`, false},
		{"string", `{ "description": "New value of axillary heater.", "example": "test" }`, false},
		{"number", `{ "description": "New value of axillary heater.", "example": 1 }`, false},
		{"values", `{ "description": "If the light should flash.", "values": [ "short", "long" ]}`, false},
		{"type mismatch", `{ "description": "New value of axillary heater.", "example": null }`, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ServiceField{}
			err := s.UnmarshalJSON([]byte(tt.data))
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
