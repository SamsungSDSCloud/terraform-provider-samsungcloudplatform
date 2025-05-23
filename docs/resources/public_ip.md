---
page_title: "samsungcloudplatform_public_ip Resource - samsungcloudplatform"
subcategory: ""
description: |-
  Provides a Public IP resource.
---

# Resource: samsungcloudplatform_public_ip

Provides a Public IP resource.


## Example Usage

```terraform
data "samsungcloudplatform_region" "region" {
}

resource "samsungcloudplatform_public_ip" "ip01" {
  description = "Public IP generated from Terraform"
  region      = data.samsungcloudplatform_region.region.location
  uplink_type = "INTERNET"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `region` (String) Region name
- `uplink_type` (String) Public IP uplinkType ('INTERNET'|'DEDICATED_INTERNET'|'SHARED_GROUP'|'SECURE_INTERNET')

### Optional

- `description` (String) Description of public IP
- `tags` (Map of String)

### Read-Only

- `id` (String) The ID of this resource.
- `ipv4` (String) IP address of public IP


