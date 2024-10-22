---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "samsungcloudplatform_firewalls Data Source - scp"
subcategory: ""
description: |-
  Provides list of firewalls
---

# samsungcloudplatform_firewalls (Data Source)

Provides list of firewalls

## Example Usage

```terraform
data "samsungcloudplatform_firewalls" "my_fws1" {
}

# Find all active firewalls
data "samsungcloudplatform_firewalls" "my_fws2" {
  vpc_id = "VPC-xxxxxx"
  filter {
    name   = "state"
    values = ["ACTIVE"]
  }
}

output "output_my_scp_fw1" {
  value = data.samsungcloudplatform_firewalls.my_fws1
}

output "output_my_scp_fw2" {
  value = data.samsungcloudplatform_firewalls.my_fws2
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filter` (Block Set) (see [below for nested schema](#nestedblock--filter))
- `target_id` (String) Target firewall resource id. (e.g. Internet Gateway, NAT Gateway, Load Balancer, ...)
- `vpc_id` (String) VPC id

### Read-Only

- `firewalls` (Block List) Firewall list (see [below for nested schema](#nestedblock--firewalls))
- `id` (String) The ID of this resource.

<a id="nestedblock--filter"></a>
### Nested Schema for `filter`

Required:

- `name` (String) Filtering target name
- `values` (List of String) Filtering values. Each matching value is appended. (OR rule)

Optional:

- `use_regex` (Boolean) Enable regex match for values


<a id="nestedblock--firewalls"></a>
### Nested Schema for `firewalls`

Read-Only:

- `id` (String) Firewall id
- `name` (String) Name of firewall
- `state` (String) Firewall status
- `target_id` (String) Target firewall resource id
- `target_type` (String) Target firewall resource type
- `vpc_id` (String) VPC id

