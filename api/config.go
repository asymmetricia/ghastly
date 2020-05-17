package api

import (
	"errors"
	"fmt"
	"strings"
)

// https://github.com/home-assistant/core/blob/master/homeassistant/components/config/config_entries.py#L79
type ConfigEntry struct {
	EntryId         EntryId `json:"entry_id"`
	Domain          string  `json:"domain"`
	Title           string  `json:"title"`
	Source          string  `json:"source"`
	State           string  `json:"state"`
	ConnectionClass string  `json:"connection_class"`
	SupportsOptions bool    `json:"supports_options"`
	client          *Client
}

type EntryId string

type SystemOptions struct {
	DisableNewEntities bool `json:"disable_new_entities"`
}

type ListSystemOptionsMessage struct {
	EntryId EntryId `json:"entry_id"`
}

func (ListSystemOptionsMessage) Type() string { return "config_entries/system_options/list" }

func (c *ConfigEntry) GetSystemOptions() (*SystemOptions, error) {
	return c.client.GetSystemOptions(c.EntryId)
}

func (c *Client) GetSystemOptions(entryId EntryId) (*SystemOptions, error) {
	ret, err := c.RawWebsocketRequestAs(ListSystemOptionsMessage{entryId}, (*SystemOptions)(nil))
	if err != nil {
		return nil, err
	}
	return ret.(*SystemOptions), nil
}

// ListConfigEntries lists known config entries. A config entry is basically a top-level device category; examples are
// `zwave` or `wemo`.
func (c *Client) ListConfigEntries() ([]*ConfigEntry, error) {
	ret, err := c.RawRESTGetAs("config/config_entries/entry", nil, ([]*ConfigEntry)(nil))
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
	Errors                  map[string]string      `json:"errors"`
	FlowId                  ConfigFlowId           `json:"flow_id"`
	Handler                 string                 `json:"handler"`
	StepId                  string                 `json:"step_id"`
	Type                    string                 `json:"type"`

	// These fields appear in the response to SetFlow, maybe based on `Type`?
	Description string `json:"description,omitempty"`
	Result      string `json:"result,omitempty"`
	Title       string `json:"title,omitempty"`
}

type ConfigFlowDataSchema struct {
	Name     string `json:"name"`
	Required bool   `json:"required"`
	Type     string `json:"type"`
}

// GetFlow returns the in-progress-but-not-started flow with the given ID.
func (c *Client) GetFlow(id ConfigFlowId) (*ConfigFlow, error) {
	retI, err := c.RawRESTGetAs("config/config_entries/flow/"+string(id), nil, (*ConfigFlow)(nil))
	if err != nil {
		return nil, err
	}

	ret, ok := retI.(*ConfigFlow)
	if !ok {
		return nil, fmt.Errorf("returned object was %T, not *ConfigFlow", retI)
	}

	return ret, nil
}

// GetOptionsFlow gets the current status of the options flow with the given ID.
func (c *Client) GetOptionsFlow(id ConfigFlowId) (*ConfigFlow, error) {
	obj, err := c.RawRESTGetAs("config/config_entries/options/flow/"+string(id), nil, (*ConfigFlow)(nil))
	if err != nil {
		return nil, err
	}
	return obj.(*ConfigFlow), nil
}

// SetFlow sets the configuration for the given FlowId to the given payload, and returns the result ID. I've observed
// result ID to refer to an Entry, but I guess it could refer to a second-stage config flow?
//
// If anything goes wrong, result will be the zero value and err will be non-nil.
func (c *Client) SetFlow(id ConfigFlowId, payload map[string]interface{}) (result string, err error) {
	obj, err := process(c.Post("config/config_entries/flow/"+string(id), payload))
	flowI, err := convert(obj, (*ConfigFlow)(nil), err)
	if err != nil {
		return "", fmt.Errorf("POSTing set-flow: %v", err)
	}
	flow, ok := flowI.(*ConfigFlow)
	if !ok {
		return "", fmt.Errorf("response object was %T, not *ConfigFlow", flowI)
	}

	var errs []string
	for k, flowErr := range flow.Errors {
		errs = append(errs, k+": "+flowErr)
	}
	if len(errs) > 0 {
		return "", fmt.Errorf("server responsed with error: %s", strings.Join(errs, "; "))
	}

	if flow.Result != "" {
		return flow.Result, nil
	}
	return "", errors.New("flow Result was empty")
}

// StartOptionsFlow initiates a options flow with the given handler, usually (always?) a ConfigEntry id. The new flow
//is returned. The UI calls DELETE if a thus-created flow is not used. It's not clear to me what happens if you don't do
// this.
func (c *Client) StartOptionsFlow(entryId EntryId) (*ConfigFlow, error) {
	obj, err := process(c.Post("config/config_entries/options/flow", map[string]string{"handler": string(entryId)}))
	flowI, err := convert(obj, (*ConfigFlow)(nil), err)
	if err != nil {
		return nil, fmt.Errorf("POSTing start-flow: %v", err)
	}

	flow, ok := flowI.(*ConfigFlow)
	if !ok {
		return nil, fmt.Errorf("response object was %T, not *ConfigFlow", flowI)
	}

	var errs []string
	for k, flowErr := range flow.Errors {
		errs = append(errs, k+": "+flowErr)
	}
	if len(errs) > 0 {
		return nil, fmt.Errorf("server responsed with error: %s", strings.Join(errs, "; "))
	}

	return flow, nil
}

func (c *Client) ListFlowHandlers() ([]string, error) {
	handlers, err := c.RawRESTGetAs("config/config_entries/flow_handlers", nil, ([]string)(nil))
	if err != nil {
		return nil, err
	}
	return handlers.([]string), nil
}
