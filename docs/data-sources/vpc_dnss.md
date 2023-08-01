---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "scp_vpc_dnss Data Source - scp"
subcategory: ""
description: |-
  Provides list of vpc DNS's.
---

# scp_vpc_dnss (Data Source)

Provides list of vpc DNS's.

## Example Usage

```terraform
data "scp_vpcs" "vpcs" {
}

data "scp_vpc_dnss" "dnss" {
  vpc_id = data.scp_vpcs.vpcs.contents[0].vpc_id
}

output "contents" {
  value = data.scp_vpc_dnss.dnss.contents
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `vpc_id` (String) VPC id

### Read-Only

- `contents` (Block List) VPC DNS list (see [below for nested schema](#nestedblock--contents))
- `id` (String) The ID of this resource.
- `total_count` (Number) Total list size

<a id="nestedblock--contents"></a>
### Nested Schema for `contents`

Read-Only:

- `dns_user_zone_domain` (String) Zone Domain
- `dns_user_zone_id` (String) Zone Domain Id
- `dns_user_zone_name` (String) Zone Name
- `dns_user_zone_server_ip` (String) Zone Dns IP
- `dns_user_zone_source_ip` (String) Zone Source IP
- `dns_user_zone_state` (String) Zone State

