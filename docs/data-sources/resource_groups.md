---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "samsungcloudplatform_resource_groups Data Source - samsungcloudplatform"
subcategory: ""
description: |-
  
---

# samsungcloudplatform_resource_groups (Data Source)



## Example Usage

```terraform
data "samsungcloudplatform_resource_groups" "my_resource_groups" {
}

output "result_my_resource_groups" {
  value = data.samsungcloudplatform_resource_groups.my_resource_groups
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `created_by_id` (String) The user id which created the resource group
- `modified_by_email` (String) The user email which modified the resource group
- `modified_by_id` (String) The user id which modified the resource group
- `resource_group_name` (String) Resource group name

### Read-Only

- `contents` (Block List) Resource group list (see [below for nested schema](#nestedblock--contents))
- `id` (String) The ID of this resource.
- `total_count` (Number) total count

<a id="nestedblock--contents"></a>
### Nested Schema for `contents`

Read-Only:

- `created_by_email` (String) The user email which created the resource group
- `created_by_id` (String) The user id which created the resource group
- `created_by_name` (String) The user name which created the resource group
- `created_dt` (String) The created date of the resource group
- `modified_by_email` (String) The user email which modified the resource group
- `modified_by_id` (String) The user id which modified the resource group
- `modified_by_name` (String) The user name which modified the resource group
- `modified_dt` (String) The modified date of the resource group
- `resource_group_description` (String) Resource group description
- `resource_group_id` (String) Resource group id
- `resource_group_name` (String) Resource group name


