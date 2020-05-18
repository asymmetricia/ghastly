# `homeassistant_entity_ids` Data Source

This data source retrieves a list of entity IDs for all entities that match the configuration. See [`entity.md`](entity.md) for an option that returns full information for entities, albeit only for one entity at a time.

## Example Usage

```hcl
data "homeassistant_entity_ids" "examples" {
  entity_id = "switch.acme_example_123"
}
```

```hcl
data "homeassistant_entity_ids" "all_switches" {
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

In addition to the above arguments, one attribute is exported:

* `entity_ids` -- a `set(string)` of discovered entity IDs.
