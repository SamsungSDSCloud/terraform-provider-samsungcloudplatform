---
page_title: "scp_direct_connect Resource - scp"
subcategory: ""
description: |-
  Provides a DirectConnect resource.
---

# Resource: scp_direct_connect

Provides a DirectConnect resource.


## Example Usage

```terraform
data "scp_region" "my_region" {
}

resource "scp_direct_connect" "dc01" {
  name        = var.name
  description = "DirectConnect generated from Terraform"
  region      = data.scp_region.my_region.location
  bandwidth   = var.bandwidth
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `bandwidth` (Number) Bandwidth gbps. (1 or 10)
- `name` (String) DirectConnect name. (3 to 20 characters without specials)
- `region` (String) Region name

### Optional

- `description` (String) DirectConnect description. (Up to 50 characters)

### Read-Only

- `id` (String) The ID of this resource.
