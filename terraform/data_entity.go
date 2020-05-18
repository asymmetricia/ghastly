package main

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pdbogen/ghastly/api"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

func entityFilter() map[string]*schema.Schema {
	entitySchema := map[string]*schema.Schema{
		// bonus input
		"entity_id_prefix": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "match entities whose entity_id has this prefix",
		},
	}

	typ := reflect.TypeOf(api.Entity{})
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		jsonName := strings.Split(field.Tag.Get("json"), ",")[0]
		entitySchema[jsonName] = &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: fmt.Sprintf("match against the entity %q field", jsonName),
		}
	}

	return entitySchema
}

func getMatchingEntities(data *schema.ResourceData, i interface{}) ([]api.Entity, error) {
	client := i.(*api.Client)
	entities, err := client.ListEntities()
	if err != nil {
		return nil, fmt.Errorf("listing entities: %w", err)
	}

	typ := reflect.TypeOf(api.Entity{})

	var filtered []api.Entity
entities:
	for _, entity := range entities {
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			attrName := strings.Split(field.Tag.Get("json"), ",")[0]
			attr, _ := data.Get(attrName).(string)
			if len(attr) == 0 {
				continue
			}

			actual := reflect.ValueOf(entity).Field(i).String()
			if actual != attr {
				continue entities
			}
		}

		entityIdPrefix, _ := data.Get("entity_id_prefix").(string)
		if entityIdPrefix != "" && !strings.HasPrefix(entity.EntityId, entityIdPrefix) {
			continue
		}

		filtered = append(filtered, entity)
	}

	return filtered, nil
}

func dataEntity() *schema.Resource {
	entitySchema := entityFilter()

	return &schema.Resource{
		Schema: entitySchema,
		Read: func(data *schema.ResourceData, i interface{}) error {
			entities, err := getMatchingEntities(data, i)
			if err != nil {
				return err
			}

			if len(entities) < 1 {
				return errors.New("no match")
			}

			if len(entities) > 1 {
				return errors.New("too many matches")
			}

			val := reflect.ValueOf(entities[0])
			typ := val.Type()
			for i := 0; i < typ.NumField(); i++ {
				data.Set(
					strings.Split(typ.Field(i).Tag.Get("json"), ",")[0],
					val.Field(i).String(),
				)
			}

			data.SetId(entities[0].EntityId)
			return nil
		},
	}
}
