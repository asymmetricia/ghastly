package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pdbogen/ghastly/api"
)

func resourceEntityName() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"entity_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
		Create: func(data *schema.ResourceData, i interface{}) error {
			client := i.(*api.Client)
			entity, err := client.GetEntity(data.Get("entity_id").(string))
			if err != nil {
				return err
			}

			if err := client.SetEntityName(entity.EntityId, data.Get("name").(string)); err != nil {
				return err
			}

			data.SetId(entity.EntityId)
			data.Set("name", data.Get("name"))
			return nil
		},
		Read: func(data *schema.ResourceData, i interface{}) error {
			client := i.(*api.Client)
			entity, err := client.GetEntity(data.Get("entity_id").(string))
			if err != nil {
				return err
			}
			data.SetId(entity.EntityId)
			data.Set("entity_id", entity.EntityId)
			data.Set("name", entity.Name)
			return nil
		},
		Delete: func(data *schema.ResourceData, i interface{}) error {
			client := i.(*api.Client)
			entity, err := client.GetEntity(data.Get("entity_id").(string))
			if err != nil {
				return err
			}

			return client.SetEntityName(entity.EntityId, "")
		},
	}
}
