package main

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pdbogen/ghastly/api"
	"strconv"
	"strings"
)

func dataEntity() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"entity_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"entity_id_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"config_entry_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"device_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"disabled_by": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"platform": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
		Read: func(data *schema.ResourceData, i interface{}) error {
			client := i.(*api.Client)
			entities, err := client.ListEntities()
			if err != nil {
				return fmt.Errorf("getting eneity list: %w", err)
			}

			var filteredEntities []api.Entity

			for _, entity := range entities {
				if eid := data.Get("entity_id").(string); eid != "" && entity.EntityId != eid {
					continue
				}
				if eidPre := data.Get("entity_id_prefix").(string); eidPre != "" && !strings.HasPrefix(entity.EntityId, eidPre) {
					continue
				}
				if confId := data.Get("config_entry_id").(string); confId != "" && entity.ConfigEntryId != confId {
					continue
				}
				if devId := data.Get("device_id").(string); devId != "" && entity.DeviceId != devId {
					continue
				}
				if disBy := data.Get("disabled_by").(string); disBy != "" && entity.DisabledBy != disBy {
					continue
				}
				if name := data.Get("name").(string); name != "" && entity.Name != name {
					continue
				}
				if platform := data.Get("platform").(string); platform != "" && entity.Platform != platform {
					continue
				}
				filteredEntities = append(filteredEntities, entity)
			}

			if len(filteredEntities) == 0 {
				return errors.New("no match")
			}
			if len(filteredEntities) > 1 {
				return errors.New("too many matches; consider narrowing parameters")
			}

			entity := filteredEntities[0]

			data.SetId(strconv.Itoa(hashcode.String(entity.EntityId)))
			err = data.Set("entity_id", entity.EntityId)
			if err == nil {
				err = data.Set("config_entry_id", entity.ConfigEntryId)
			}
			if err == nil {
				err = data.Set("device_id", entity.DeviceId)
			}
			if err == nil {
				err = data.Set("disabled_by", entity.DisabledBy)
			}
			if err == nil {
				err = data.Set("name", entity.Name)
			}
			if err == nil {
				err = data.Set("platform", entity.Platform)
			}
			return err
		},
	}
}
