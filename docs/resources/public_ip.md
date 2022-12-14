---
page_title: "scp_public_ip Resource - scp"
subcategory: ""
description: |-
  Provides a Public IP resource.
---

# Resource: scp_public_ip

Provides a Public IP resource.


## Example Usage

```terraform
data "scp_region" "region" {
}

resource "scp_public_ip" "ip01" {
  description = "Public IP generated from Terraform"
  region      = data.scp_region.region.location
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `region` (String) Region name

### Optional

- `description` (String) Description of public IP

### Read-Only

- `id` (String) The ID of this resource.
- `ipv4` (String) IP address of public IP
