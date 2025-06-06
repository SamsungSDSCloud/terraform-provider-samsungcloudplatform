---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "samsungcloudplatform_project_zones Data Source - samsungcloudplatform"
subcategory: ""
description: |-
  Provides list of service zones in project
---

# samsungcloudplatform_project_zones (Data Source)

Provides list of service zones in project

## Example Usage

```terraform
data "samsungcloudplatform_project_zones" "my_scp_project_zones" {
  project_id = "PROJECT-XXXXXXX"
}

output "output_my_scp_project" {
  value = data.samsungcloudplatform_project_zones.my_scp_project_zones
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `project_id` (String) Project ID

### Optional

- `filter` (Block Set) (see [below for nested schema](#nestedblock--filter))

### Read-Only

- `contents` (Block List) Zones in project (see [below for nested schema](#nestedblock--contents))
- `id` (String) The ID of this resource.
- `total_count` (Number) Total list size

<a id="nestedblock--filter"></a>
### Nested Schema for `filter`

Required:

- `name` (String) Filtering target name
- `values` (List of String) Filtering values. Each matching value is appended. (OR rule)

Optional:

- `use_regex` (Boolean) Enable regex match for values


<a id="nestedblock--contents"></a>
### Nested Schema for `contents`

Read-Only:

- `availability_zones` (Block List) List of availability zones (see [below for nested schema](#nestedblock--contents--availability_zones))
- `block_id` (String) Block ID
- `is_multi_availability_zone` (Boolean) Multi availability zone
- `service_zone_id` (String) Service zone ID
- `service_zone_location` (String) Service zone location
- `service_zone_name` (String) Service zone name

<a id="nestedblock--contents--availability_zones"></a>
### Nested Schema for `contents.availability_zones`

Read-Only:

- `availability_zone_name` (String) Availability zone name


