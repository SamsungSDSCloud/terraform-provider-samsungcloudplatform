---
page_title: "samsungcloudplatform_dns_domain Resource - samsungcloudplatform"
subcategory: ""
description: |-
  Provides a Dns Domain resource. (Only available for PRIVATE environment usage type)
---

# Resource: samsungcloudplatform_dns_domain

Provides a Dns Domain resource. (Only available for PRIVATE environment usage type)


## Example Usage

```terraform
resource "samsungcloudplatform_dns_domain" "my_dns_domain" {
  dns_domain_name       = var.name
  dns_root_domain_name  = var.root_domain_name
  dns_description       = "terraform test 2"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `dns_domain_name` (String) DNS Name
- `dns_root_domain_name` (String) DNS Root Domain Name

### Optional

- `dns_description` (String) DNS Domain Description
- `tags` (Map of String)
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `update` (String)


