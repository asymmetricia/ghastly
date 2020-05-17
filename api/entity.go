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
	ConfigEntryId string  `json:"config_entry_id,omitempty"`
	DeviceId      string  `json:"device_id,omitempty"`
	DisabledBy    *string `json:"disabled_by,omitempty"`
	EntityId      string  `json:"entity_id,omitempty"`
	Platform      string  `json:"platform,omitempty"`
	Name          string  `json:"name,omitempty"`
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

type EntityList struct{}

func (c *Client) ListEntities() ([]Entity, error) {
	entsI, err := c.RawWebsocketRequestAs(EntityListMessage{}, []Entity{})
	if err != nil {
		return nil, err
	}

	entities, ok := entsI.([]Entity)
	if !ok {
		return nil, fmt.Errorf("got back %T, not []Entity", entsI)
	}

	return entities, nil
}
