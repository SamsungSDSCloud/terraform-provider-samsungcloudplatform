---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "samsungcloudplatform_launch_configuration Data Source - samsungcloudplatform"
subcategory: ""
description: |-
  Provides details of Launch Configuration
---

# samsungcloudplatform_launch_configuration (Data Source)

Provides details of Launch Configuration

## Example Usage

```terraform
# Find details of Launch Configuration
data "samsungcloudplatform_launch_configuration" "my_launch_configuration1" {
  lc_id = "LAUNCH_CONFIGURATION-XXXXX"
}

output "output_scp_launch_configuration1" {
  value = data.samsungcloudplatform_launch_configuration.my_launch_configuration1
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `lc_id` (String) Launch Configuration ID

### Optional

- `asg_ids` (List of String) Auto-Scaling Group ID list

### Read-Only

- `block_id` (String) Block ID
- `block_storages` (Block List) Block Storage list (see [below for nested schema](#nestedblock--block_storages))
- `contract_product_id` (String) Contract product ID
- `created_by` (String) The person who created the resource
- `created_dt` (String) Creation date
- `id` (String) The ID of this resource.
- `image_id` (String) Image ID
- `initial_script` (String) Virtual Server's initial script
- `key_pair_id` (String) Key pair ID
- `lc_name` (String) Launch Configuration name
- `modified_by` (String) The person who modified the resource
- `modified_dt` (String) Modification date
- `os_product_id` (String) OS product ID
- `os_type` (String) OS type
- `product_group_id` (String) Product group ID
- `project_id` (String) Project ID
- `scale_product_id` (String) Scale product ID
- `server_type` (String) Server type
- `service_zone_id` (String) Service zone ID

<a id="nestedblock--block_storages"></a>
### Nested Schema for `block_storages`

Read-Only:

- `block_storage_size` (Number) Block Storage size (GB)
- `disk_type` (String) Block storage product (default value : SSD)
- `encryption_enabled` (Boolean) Encryption enabled
- `is_boot_disk` (Boolean) Is boot disk or not
- `product_id` (String) Product ID


