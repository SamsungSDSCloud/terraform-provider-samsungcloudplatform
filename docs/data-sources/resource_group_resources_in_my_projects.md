---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "samsungcloudplatform_resource_group_resources_in_my_projects Data Source - samsungcloudplatform"
subcategory: ""
description: |-
  
---

# samsungcloudplatform_resource_group_resources_in_my_projects (Data Source)



## Example Usage

```terraform
data "samsungcloudplatform_resource_group_resources_in_my_projects" "my_resource_group_resources_in_my_projects" {
  resource_group_id = "RESOURCE_GROUP-XXXXXXXXXXXXX"
}

output "result_my_resource_group_resources" {
  value = data.samsungcloudplatform_resource_group_resources_in_my_projects.my_resource_group_resources_in_my_projects
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `resource_group_id` (String) Resource group id

### Optional

- `contents` (List of Map of String) Resource list
- `created_by_id` (String) The user id which created the resource
- `modified_by_id` (String) The user id which modified the resource
- `resource_id` (String) Resource id
- `resource_name` (String) Resource name

### Read-Only

- `id` (String) The ID of this resource.
- `total_count` (Number) total count


