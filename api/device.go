package api

import (
	"encoding/json"
	"errors"
	"fmt"
)

type DeviceListMessage struct{}

func (DeviceListMessage) Type() string { return "config/device_registry/list" }

func init() {
	RegisterMessageType(DeviceListMessage{})
}

type Device struct {
	ID            string      `json:"id"`
	AreaId        *string     `json:"area_id"`
	ConfigEntries []string    `json:"config_entries"`
	Connections   [][2]string `json:"connections"`
	Manufacturer  *string     `json:"manufacturer"`
	Model         *string     `json:"model"`
	Name          *string     `json:"name"`
	NameByUser    *string     `json:"name_by_user"`
	SwVersion     *string     `json:"sw_version"`
	ViaDeviceId   *string     `json:"via_device_id"`
}

func (c *Client) GetDevice(id string) (*Device, error) {
	devices, err := c.ListDevices()
	if err != nil {
		return nil, err
	}

	for _, device := range devices {
		if device.ID == id {
			return device, nil
		}
	}
	return nil, errors.New("not found")
}

func (c *Client) ListDevices() ([]*Device, error) {
	retI, err := c.Exchange(DeviceListMessage{})
	if err != nil {
		return nil, err
	}
	result, ok := retI.(*ResultMessage)
	if !ok {
		return nil, fmt.Errorf("server sent %T, not result", retI)
	}
	if !result.Success {
		return nil, fmt.Errorf("get unsuccessful: %v", result.Error)
	}

	resultsJson, _ := json.Marshal(result.Result)
	var ret []*Device
	if err := json.Unmarshal(resultsJson, &ret); err != nil {
		return nil, fmt.Errorf("converting results: %v", err)
	}

	return ret, nil
}
