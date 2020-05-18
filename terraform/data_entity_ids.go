package main

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pdbogen/ghastly/api"
	"strings"
)

func dataEntityIds() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"device_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"platform": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"entity_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
		Read: func(data *schema.ResourceData, i interface{}) error {
			client := i.(*api.Client)
			entities, err := client.ListEntities()
			if err != nil {
				return fmt.Errorf("listing entities: %w", err)
			}

			deviceId := data.Get("device_id").(string)
			platform := data.Get("platform").(string)
			prefix := data.Get("prefix").(string)

			var entityIds []string
			for _, entity := range entities {
				if len(deviceId) > 0 && entity.DeviceId != deviceId {
					continue
				}
				if len(platform) > 0 && entity.Platform != platform {
					continue
				}
				if len(prefix) > 0 && !strings.HasPrefix(entity.EntityId, prefix) {
					continue
				}
				entityIds = append(entityIds, entity.EntityId)
			}

			data.SetId(hashcode.Strings([]string{deviceId, platform, prefix}))
			data.Set("entity_ids", entityIds)

			return nil
		},
	}
}
