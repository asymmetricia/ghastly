package api

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
)

type EntityListMessage struct{}

func (EntityListMessage) Type() string { return "config/entity_registry/list" }

type EntityGetMessage struct {
	EntityId string `json:"entity_id"`
}

func (EntityGetMessage) Type() string { return "config/entity_registry/get" }

func init() {
	RegisterMessageType(EntityListMessage{})
}

type Entity struct {
	ID string
}

func (c *Client) GetEntity(id string) (*Entity, error) {
	retI, err := c.Exchange(EntityGetMessage{id})
	if err != nil {
		return nil, err
	}
	ret, ok := retI.(*ResultMessage)
	if !ok {
		return nil, fmt.Errorf("server sent %T, not result", retI)
	}
	if !ret.Success {
		return nil, fmt.Errorf("get unsuccessful: %v", ret.Error)
	}
	logrus.Debug(ret)
	return nil, errors.New("unimplemented")
}
