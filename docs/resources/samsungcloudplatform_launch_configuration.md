---
page_title: "samsungcloudplatform_launch_configuration Resource - scp"
subcategory: ""
description: |-
  Provides a Launch Configuration resource.
---

# Resource: samsungcloudplatform_launch_configuration

Provides a Launch Configuration resource.


## Example Usage

```terraform
resource "samsungcloudplatform_launch_configuration" "my_launch_configuration" {
  dynamic "block_storages" {
    for_each = var.block_storages
    content {
      block_storage_size = block_storages.value["block_storage_size"]
      disk_type = block_storages.value["disk_type"]
      encryption_enabled = block_storages.value["encryption_enabled"]
      is_boot_disk = block_storages.value["is_boot_disk"]
    }
  }
  image_id = var.image_id
  initial_script = var.initial_script
  key_pair_id = var.key_pair_id
  lc_name = var.lc_name
  server_type = var.server_type
  service_zone_id = var.service_zone_id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `block_storages` (Block List, Min: 1) Block Storage list (see [below for nested schema](#nestedblock--block_storages))
- `image_id` (String) Image ID
- `key_pair_id` (String) Key pair ID
- `lc_name` (String) Launch Configuration name
- `server_type` (String) Server type
- `service_zone_id` (String) Service zone ID

### Optional

- `initial_script` (String) Virtual Server's initial script
- `tags` (Map of String)

### Read-Only

- `asg_ids` (List of String) Auto-Scaling Group ID list
- `block_id` (String) Block ID
- `contract_product_id` (String) Contract product ID
- `created_by` (String) The person who created the resource
- `created_dt` (String) Creation date
- `id` (String) The ID of this resource.
- `lc_id` (String) Launch Configuration ID
- `modified_by` (String) The person who modified the resource
- `modified_dt` (String) Modification date
- `os_product_id` (String) OS product ID
- `os_type` (String) OS type
- `product_group_id` (String) Product group ID
- `project_id` (String) Project ID
- `scale_product_id` (String) Scale product ID

<a id="nestedblock--block_storages"></a>
### Nested Schema for `block_storages`

Required:

- `block_storage_size` (Number) Block Storage size (GB)
- `disk_type` (String) Block Storage product (default value : SSD)
- `encryption_enabled` (Boolean) Encryption enabled
- `is_boot_disk` (Boolean) Is boot disk or not

Read-Only:

- `product_id` (String) Product ID