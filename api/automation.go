package api

import (
	"fmt"
	"strings"
	"time"
)

type AutomationId string

type AutomationListEntry struct {
	FriendlyName  string       `json:"friendly_name"`
	Id            AutomationId `json:"id"`
	LastTriggered time.Time    `json:"last_triggered"`
}

type Automation struct {
	Action    AutomationAction      `json:"action"`
	Alias     string                `json:"alias"`
	Condition []AutomationCondition `json:"condition,omitempty"`
	Id        AutomationId          `json:"id"`
	Trigger   []AutomationTrigger   `json:"trigger"`
}

func (c *Client) GetAutomation(id AutomationId) (*Automation, error) {
	retI, err := c.RawRESTGetAs("config/automation/config/"+string(id), nil, (*Automation)(nil))
	if err != nil {
		return nil, err
	}

	return retI.(*Automation), nil
}

func (c *Client) ListAutomations() ([]AutomationListEntry, error) {
	states, err := c.ListStates()
	if err != nil {
		return nil, fmt.Errorf("automations are discovered via states, but could not get states: %w", err)
	}

	var ret []AutomationListEntry
	for _, state := range states {
		if !strings.HasPrefix(state.EntityId, "automation.") {
			continue
		}

		entry := AutomationListEntry{}
		entry.Id = AutomationId(state.Context.Id)
		attr, ok := state.Attributes["friendly_name"]
		if ok {
			entry.FriendlyName, _ = attr.(string)
		}
		attr, ok = state.Attributes["last_triggered"]
		if ok {
			lt, _ := attr.(string)
			entry.LastTriggered, err = time.Parse(time.RFC3339, lt)
			if err != nil {
				return nil, fmt.Errorf("could not parse last_triggered time %q: %w", lt, err)
			}
		}

		ret = append(ret, entry)
	}

	return ret, nil
}
