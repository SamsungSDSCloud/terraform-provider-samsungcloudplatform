---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "samsungcloudplatform_bm_block_storages Data Source - samsungcloudplatform"
subcategory: ""
description: |-
  Provides Block Storage(BM) List
---

# samsungcloudplatform_bm_block_storages (Data Source)

Provides Block Storage(BM) List

## Example Usage

```terraform
data "samsungcloudplatform_bm_block_storages" "my_scp_bm_block_storages" {
}

output "output_my_scp_bm_block_storages_org" {
  value = data.samsungcloudplatform_bm_block_storages.my_scp_bm_block_storages
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `contents` (Block List) BareMetal Block Storages (see [below for nested schema](#nestedblock--contents))

### Read-Only

- `id` (String) The ID of this resource.
- `total_count` (Number) Total list size

<a id="nestedblock--contents"></a>
### Nested Schema for `contents`

Optional:

- `bare_metal_server_ids` (List of String) Baremetal Server Ids

Read-Only:

- `bare_metal_block_storage_id` (String) Baremetal Block Storage Id
- `bare_metal_block_storage_name` (String) Baremetal Block Storage Name
- `bare_metal_block_storage_purpose` (String) Baremetal Block Storage Purpose
- `bare_metal_block_storage_size` (Number) Baremetal Block Storage Size
- `bare_metal_block_storage_state` (String) Baremetal Block Storage State
- `bare_metal_block_storage_type_id` (String) Baremetal Block Storage Type
- `block_id` (String) Block Id
- `created_by` (String) Created By
- `created_dt` (String) Created Date
- `encryption_enabled` (Boolean) Encryption Enabled
- `location` (String) Location
- `modified_by` (String) Modified By
- `modified_dt` (String) Modified Date
- `service_zone_id` (String) Service Zone Id


