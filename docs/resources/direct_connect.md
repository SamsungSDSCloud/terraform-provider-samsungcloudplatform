---
page_title: "samsungcloudplatform_direct_connect Resource - samsungcloudplatform"
subcategory: ""
description: |-
  Provides a DirectConnect resource.
---

# Resource: samsungcloudplatform_direct_connect

Provides a DirectConnect resource.


## Example Usage

```terraform
data "samsungcloudplatform_region" "my_region" {
}

resource "samsungcloudplatform_direct_connect" "dc01" {
  name        = var.name
  description = "DirectConnect generated from Terraform"
  region      = data.samsungcloudplatform_region.my_region.location
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
- `tags` (Map of String)

### Read-Only

- `id` (String) The ID of this resource.


