---
page_title: "samsungcloudplatform_migration_image Resource - samsungcloudplatform"
subcategory: ""
description: |-
  Provides Migration Image resource.
---

# Resource: samsungcloudplatform_migration_image

Provides Migration Image resource.


## Example Usage

```terraform
resource "samsungcloudplatform_migration_image" "my_migration_image" {
  image_name = var.name
  original_image_id = var.image_id
  ova_url = var.url
  access_key = var.access_key
  secret_key = var.secret_key
  os_user_id = var.os_id
  os_user_password = var.os_pw
  image_description = var.desc
  az_name =var.az_name
  service_zone_id = var.service_zone_id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `access_key` (String) access key for ova
- `image_description` (String) Image Description
- `image_name` (String) Migration Image Name
- `original_image_id` (String) Original Image Id
- `os_user_id` (String) OS User Id
- `os_user_password` (String) Os User Password
- `ova_url` (String) Ova url
- `secret_key` (String) secret key for ova
- `service_zone_id` (String)

### Optional

- `az_name` (String) Availability Zone Name
- `icon` (Map of String)
- `properties` (Map of String)
- `tags` (Map of String)

### Read-Only

- `created_by` (String)
- `created_dt` (String)
- `id` (String) The ID of this resource.
- `image_id` (String)
- `image_state` (String) Image state (ACTIVE)
- `image_type` (String) Image type (STANDARD, CUSTOM, MIGRATION)
- `modified_by` (String)
- `modified_dt` (String)
- `origin_image_name` (String)
- `os_type` (String) OS type (Windows, Ubuntu, ..)
- `product_group_id` (String)
- `products` (Block List) (see [below for nested schema](#nestedblock--products))

<a id="nestedblock--products"></a>
### Nested Schema for `products`

Read-Only:

- `created_dt` (String)
- `image_id` (String)
- `product_id` (String)
- `product_name` (String)
- `product_type` (String)
- `product_value` (String)
- `seq` (Number)


