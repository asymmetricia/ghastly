# `homeassistant_entity` Data Source

This data source retrieves information about an entity. It can use any of its attributes to match a specific entity, but requires the only one entity be returned. See [`entity_ids.md`](entity_ids.md) for an option to discover multiple entities.

## Example Usage

```hcl
data "homeassistant_entity" "example" {
  entity_id = "switch.acme_example_123"
}
```

```hcl
data "homeassistant_entity" "example" {
  device_id        = 2520ab9a70a574ea1f7896d0281d9274
  entity_id_prefix = "switch."
}
```

## Argument Reference

* `entity_id` -- (optional) match based on entity ID
* `entity_id_prefix` -- (optional) match entities with an entity ID with this prefix
* `config_entry_id` -- (optional) match based on config entry ID
* `device_id` -- (optional) match based on device ID
* `disabled_by` -- (optional) match based on mechanism responsible for disabling this entity
* `name` -- (optional) match based on friendly name
* `platform` -- (optional) match based on platform name

## Attribute Reference

All attributes noted above as arguments are exported.
