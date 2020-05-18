package main

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pdbogen/ghastly/api"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ConfigureFunc: func(data *schema.ResourceData) (interface{}, error) {
			client := &api.Client{
				Token:  data.Get("token").(string),
				Server: data.Get("url").(string),
			}
			_, err := client.GetConfig()
			if err != nil {
				return nil, fmt.Errorf("getting config: %w", err)
			}
			return client, nil
		},
		Schema: map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"token": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"homeassistant_entity_name": resourceEntityName(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"homeassistant_entity":     dataEntity(),
			"homeassistant_entity_ids": dataEntityIds(),
		},
	}
}
