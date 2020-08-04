package main

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataEntityIds() *schema.Resource {
	entitySchema := entityFilter()

	entitySchema["entity_ids"] = &schema.Schema{
		Type:     schema.TypeSet,
		Elem:     &schema.Schema{Type: schema.TypeString},
		Computed: true,
	}

	return &schema.Resource{
		Schema: entitySchema,
		Read: func(data *schema.ResourceData, i interface{}) error {
			entities, err := getMatchingEntities(data, i)
			if err != nil {
				return fmt.Errorf("listing entities: %w", err)
			}

			var entityIds []string
			for _, entity := range entities {
				entityIds = append(entityIds, entity.EntityId)
			}

			data.SetId(hashcode.Strings(entityIds))
			data.Set("entity_ids", entityIds)
			return nil
		},
	}
}
