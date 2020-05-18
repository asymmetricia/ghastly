# `homeassistant_entity_name` Resource

This resource sets the friendly name of an entity, given its entity ID.

## Example Usage

```hcl
resource "homeassistant_entity_name" "example" {
  entity_id = "switch.acme_inc_example_123"
  name      = "Example Switch"
}
```

## Argument Reference

* `entity_id` - (Required) The entity ID to have its friendly name changed.
* `name` - (Required) The friendly name to give the device

## Attribute Reference

No additional attributes are exported.
