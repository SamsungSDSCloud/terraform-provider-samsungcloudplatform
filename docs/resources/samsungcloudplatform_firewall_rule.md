---
page_title: "samsungcloudplatform_firewall_rule Resource - scp"
subcategory: ""
description: |-
  Provides a Firewall Rule resource.
---

# Resource: samsungcloudplatform_firewall_rule

Provides a Firewall Rule resource.


## Example Usage

```terraform
data "samsungcloudplatform_region" "region" {
}

resource "samsungcloudplatform_firewall_rule" "vpc4fw_fwrule" {
  firewall_id = "FIREWALL-xxxxxxx"

  direction = "IN_OUT"
  action    = "ALLOW"

  enabled = false

  source_addresses_ipv4      = ["128.0.0.0/1"]
  destination_addresses_ipv4 = ["128.0.0.0/1"]

  service {
    type  = "TCP"
    value = "8080"
  }
  service {
    type  = "UDP"
    value = "22"
  }
  service {
    type = "TCP_ALL"
    value = ""
  }

  description = "Rule from terraform"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `action` (String) Rule action. (ALLOW, DROP)
- `destination_addresses_ipv4` (List of String) Destination ip addresses list
- `direction` (String) Rule direction. (IN, OUT, IN_OUT)
- `firewall_id` (String) Firewall id
- `service` (Block Set, Min: 1) Firewall Rule service (see [below for nested schema](#nestedblock--service))
- `source_addresses_ipv4` (List of String) Source ip addresses list

### Optional

- `description` (String) Rule description. (0 to 100 characters)
- `enabled` (Boolean) Rule enabled state.
- `location_rule_id` (String) Location Rule id
- `rule_location_type` (String) Rule location type. (FIRST, BEFORE, AFTER, LAST)

### Read-Only

- `id` (String) The ID of this resource.
- `target_id` (String) Target firewall resource id

<a id="nestedblock--service"></a>
### Nested Schema for `service`

Required:

- `type` (String) Protocol type. (TCP, UDP, ICMP, ALL)

Optional:

- `value` (String) Port value