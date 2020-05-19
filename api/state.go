package api

import (
	"fmt"
	"time"
)

type State struct {
	Attributes map[string]interface{} `json:"attributes"`
	Context    struct {
		Id       string  `json:"id"`
		ParentId *string `json:"parent_id"`
		UserId   *string `json:"user_id"`
	} `json:"context"`
	EntityId    string    `json:"entity_id"`
	LastChanged time.Time `json:"last_Changed"`
	LastUpdated time.Time `json:"last_updated"`
	State       string    `json:"state"`
}

type ListStatesMessage struct{}

func (g ListStatesMessage) Type() string {
	return "get_states"
}

func (c *Client) ListStates() ([]State, error) {
	retI, err := c.RawWebsocketRequestAs(ListStatesMessage{}, ([]State)(nil))
	if err != nil {
		return nil, err
	}

	ret, ok := retI.([]State)
	if !ok {
		return nil, fmt.Errorf("received %T instead of []State", retI)
	}

	return ret, nil
}
