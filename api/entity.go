package api

import (
	"fmt"
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
	ConfigEntryId string `json:"config_entry_id,omitempty"`
	DeviceId      string `json:"device_id,omitempty"`
	DisabledBy    string `json:"disabled_by,omitempty"`
	EntityId      string `json:"entity_id,omitempty"`
	Platform      string `json:"platform,omitempty"`
	Name          string `json:"name,omitempty"`
}

func (c *Client) GetEntity(id string) (*Entity, error) {
	entityI, err := c.RawWebsocketRequestAs(EntityGetMessage{id}, (*Entity)(nil))
	if err != nil {
		return nil, err
	}

	entity, ok := entityI.(*Entity)
	if !ok {
		return nil, fmt.Errorf("server sent %T, not *Entity", entityI)
	}

	return entity, nil
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

type EntityRename struct {
	EntityId string `json:"entity_id"`
	Name     string `json:"name"`
}

func (EntityRename) Type() string { return "config/entity_registry/update" }

func (c *Client) SetEntityName(id string, name string) error {
	_, err := c.RawWebsocketRequest(EntityRename{EntityId: id, Name: name})
	return err
}
