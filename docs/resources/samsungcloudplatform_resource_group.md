---
page_title: "samsungcloudplatform_resource_group Resource - scp"
subcategory: ""
description: |-
  
---

# Resource: samsungcloudplatform_resource_group




## Example Usage

```terraform
resource "samsungcloudplatform_resource_group" "my_resource_group" {
  name = var.name

  target_resource_tags {
    tag_key = "tk01"
    tag_value = "tv01"
  }
  target_resource_tags {
    tag_key = "tk02"
    tag_value = "tv02"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Resource group name

### Optional

- `created_by_id` (String) The user id which created the resource group
- `created_by_name` (String) The user name which created the resource group
- `modified_by_id` (String) The user id which modified the resource group
- `modified_by_name` (String) The user name which modified the resource group
- `resource_group_description` (String) Resource group description
- `target_resource_tags` (Map of String)
- `target_resource_types` (List of String) Resource group types

### Read-Only

- `created_by_email` (String) The user email which created the resource group
- `created_dt` (String) The created date of the resource group
- `id` (String) The ID of this resource.
- `modified_by_email` (String) The user email which modified the resource group
- `modified_dt` (String) The modified date of the resource group
- `resource_group_name` (String) Resource group name
- `target_resource_tag` (List of Object) Tag list (see [below for nested schema](#nestedatt--target_resource_tag))

<a id="nestedatt--target_resource_tag"></a>
### Nested Schema for `target_resource_tag`

Read-Only:

- `tag_key` (String)
- `tag_value` (String)