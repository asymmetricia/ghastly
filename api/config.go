package api

import (
	"errors"
)

// https://github.com/home-assistant/core/blob/master/homeassistant/components/config/config_entries.py#L79
type ConfigEntry struct {
	EntryId         string `json:"entry_id"`
	Domain          string `json:"domain"`
	Title           string `json:"title"`
	Source          string `json:"source"`
	State           string `json:"state"`
	ConnectionClass string `json:"connection_class"`
	SupportsOptions bool   `json:"supports_options"`
	client          *Client
}

type SystemOptions struct {
	DisableNewEntities bool `json:"disable_new_entities"`
}

type ListSystemOptionsMessage struct {
	EntryId string `json:"entry_id"`
}

func (ListSystemOptionsMessage) Type() string { return "config_entries/system_options/list" }

func (c *ConfigEntry) GetSystemOptions() (*SystemOptions, error) {
	return c.client.GetSystemOptions(c.EntryId)
}

func (c *Client) GetSystemOptions(entryId string) (*SystemOptions, error) {
	ret, err := c.RawWebsocketRequestAs(ListSystemOptionsMessage{entryId}, (*SystemOptions)(nil))
	if err != nil {
		return nil, err
	}
	return ret.(*SystemOptions), nil
}

// ListConfigEntries lists known config entries. A config entry is basically a top-level device category; examples are
// `zwave` or `wemo`.
func (c *Client) ListConfigEntries() ([]*ConfigEntry, error) {
	ret, err := c.RawRESTRequestAs("config/config_entries/entry", nil, ([]*ConfigEntry)(nil))
	if err != nil {
		return nil, err
	}

	for _, r := range ret.([]*ConfigEntry) {
		r.client = c
	}

	return ret.([]*ConfigEntry), nil
}

type ConfigFlow struct{}

// ListConfigFlows would list known config flows, but the requisite API endpoint is not implemented.
func (c *Client) ListConfigFlows() ([]*ConfigFlow, error) {
	return nil, errors.New("unimplemented")
}
