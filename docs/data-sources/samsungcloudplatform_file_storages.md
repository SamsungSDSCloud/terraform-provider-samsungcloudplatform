---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "samsungcloudplatform_file_storages Data Source - scp"
subcategory: ""
description: |-
  Provides list of file storages
---

# samsungcloudplatform_file_storages (Data Source)

Provides list of file storages

## Example Usage

```terraform
data "samsungcloudplatform_file_storages" "my_scp_file_storages" {
  file_storage_states = [
    "ACTIVE",
    "ERROR"
  ]
  sort = [
    "fileStorageName:DESC"
  ]
}

output "output_my_scp_file_storages" {
  value = data.samsungcloudplatform_file_storages.my_scp_file_storages
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `block_id` (String) Block ID
- `created_by` (String) Created By
- `file_storage_id` (String) File Storage ID
- `file_storage_name` (String) File Storage Name
- `file_storage_protocol` (String) File Storage Protocol
- `file_storage_state` (String) File Storage State
- `file_storage_states` (List of String) File Storage States
- `page` (Number) Page start number from which to get the list
- `service_zone_id` (String) Service Zone ID
- `size` (Number) Size to get list
- `sort` (List of String) Sort

### Read-Only

- `contents` (List of Object) File Storage List (see [below for nested schema](#nestedatt--contents))
- `id` (String) The ID of this resource.
- `total_count` (Number) Total List Size

<a id="nestedatt--contents"></a>
### Nested Schema for `contents`

Read-Only:

- `block_id` (String)
- `created_by` (String)
- `created_dt` (String)
- `disk_type` (String)
- `encryption_enabled` (Boolean)
- `file_storage_id` (String)
- `file_storage_name` (String)
- `file_storage_protocol` (String)
- `file_storage_purpose` (String)
- `file_storage_state` (String)
- `linked_object_count` (Number)
- `modified_by` (String)
- `modified_dt` (String)
- `product_group_id` (String)
- `project_id` (String)
- `service_zone_id` (String)
- `tiering_enabled` (Boolean)

