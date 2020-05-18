# HomeAssistant Provider

Summary of what the provider is for, including use cases and links to
app/service documentation.

## Example Usage

### Change the friendly name of a specific entity

```hcl
resource "homeassistant_entity_name" "example" {
  entity_id = "switch.acme_inc_example_123"
  name = "Example Switch"
}
```

### Change the friendly names of all entities tied to a device

```hcl
data "homeassistant_entity_ids" "example_ids" {
  device_id = "9e52f20df740bff33dd9c06c51189ea6"
}

// this is our root device, we use its name to figure out how to clean up friendly names
data "homeassistant_entity" "example_base" {
  device_id        = "9e52f20df740bff33dd9c06c51189ea6"
  entity_id_prefix = "zwave."
}

locals {
  base = split(".", data.homeassistant_entity.example_base.entity_id)[1]
}

// this will let us get the name of each device to feed into the rename resource
data "homeassistant_entity" "example_all" {
  for_each  = data.homeassistant_entity_ids.example_ids.entity_ids
  entity_id = each.value
}

resource "homeassistant_entity_name" "example" {
  for_each  = data.homeassistant_entity_ids.example_ids.entity_ids
  entity_id = each.value
  name      = "example ${replace(trimprefix(split(".", data.homeassistant_entity.example_all[each.value].entity_id)[1], local.base), "_", " " )}"
}

//  name      = "example ${                                              // 5. Pre-pend our chosen friendly name
//    replace(                                                           // 4. Replace any _ with space
//      trimprefix(                                                      // 3. Remove the common base prefix
//        split(                                                         // 2. Remove the "device type" prefix
//          ".",
//          data.homeassistant_entity.example_all[each.value].entity_id  // 1. Start with the entity ID
//        )[1],
//        local.base
//      ),
//      "_",
//      " "
//    )
//  }"
```

## Argument Reference

HomeAssistant provider requires both of two arguments:

* `url` -- The URL to the homeassistant server. E.g., `https://example.com:1234`
* `token` -- A HomeAssistant API token with which to authenticate
