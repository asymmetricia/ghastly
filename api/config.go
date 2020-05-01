package api

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type Config struct {
	Components            []string          `json:"components"`
	ConfigDir             string            `json:"config_dir"`
	ConfigSource          string            `json:"config_source"`
	Elevation             int               `json:"elevation"`
	Latitude              float64           `json:"latitude"`
	LocationName          string            `json:"location_name"`
	Longitude             float64           `json:"longitude"`
	TimeZone              string            `json:"time_zone"`
	UnitSystem            map[string]string `json:"unit_system"`
	Version               string            `json:"version"`
	WhitelistExternalDirs []string          `json:"whitelist_external_dirs"`
}

type GetConfigMessage struct{}

func (GetConfigMessage) Type() string { return "get_config" }

func (c *Client) GetConfig() (*Config, error) {
	retI, err := c.RawWebsocketRequestAs(GetConfigMessage{}, &Config{})
	if err != nil {
		return nil, err
	}

	ret, ok := retI.(*Config)
	if !ok {
		return nil, fmt.Errorf("returned object was %T, not *Config", retI)
	}

	return ret, nil
}

type ConfigFlowProgress struct {
	Context ConfigFlowProgressContext `json:"context"`
	FlowId  ConfigFlowId              `json:"flow_id"`
	Handler string                    `json:"handler"`
}

type ConfigFlowProgressContext struct {
	Source string `json:"source"`
}

type ConfigFlowId string

type ListConfigFlowProgressMessage struct{}

func (ListConfigFlowProgressMessage) Type() string { return "config_entries/flow/progress" }

func (c *Client) ListConfigFlowProgress() ([]ConfigFlowProgress, error) {
	retI, err := c.RawWebsocketRequestAs(ListConfigFlowProgressMessage{}, []ConfigFlowProgress{})
	if err != nil {
		return nil, err
	}

	ret, ok := retI.([]ConfigFlowProgress)
	if !ok {
		return nil, fmt.Errorf("returned object was %T, not []ConfigFlowProgress", retI)
	}

	return ret, nil
}

type ConfigFlow struct {
	DataSchema              []ConfigFlowDataSchema `json:"data_schema"`
	DescriptionPlaceholders interface{}            `json:"description_placeholders"`
	Errors                  map[string]interface{} `json:"errors"`
	FlowId                  ConfigFlowId           `json:"flow_id"`
	Handler                 string                 `json:"handler"`
	StepId                  string                 `json:"step_id"`
	Type                    string                 `json:"type"`
}

type ConfigFlowDataSchema struct {
	Name     string `json:"name"`
	Required bool   `json:"required"`
	Type     string `json:"type"`
}

func (c *Client) GetFlow(id ConfigFlowId) (*ConfigFlow, error) {
	retI, err := c.RawRESTRequestAs("config/config_entries/flow/"+string(id), nil, (*ConfigFlow)(nil))
	if err != nil {
		return nil, err
	}

	ret, ok := retI.(*ConfigFlow)
	if !ok {
		return nil, fmt.Errorf("returned object was %T, not *ConfigFlow")
	}

	return ret, nil
}

func (c *Client) SetFlow(id ConfigFlowId, payload map[string]interface{}) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("could not marshal payload to JSON: %v", err)
	}

	_, err = process(c.Post("config/config_entries/flow/"+string(id), bytes.NewBuffer(jsonPayload)))
	if err != nil {
		return fmt.Errorf("POSTing set-flow: %v", err)
	}

	return nil
}
