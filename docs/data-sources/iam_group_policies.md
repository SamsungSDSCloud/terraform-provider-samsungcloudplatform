---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "samsungcloudplatform_iam_group_policies Data Source - samsungcloudplatform"
subcategory: ""
description: |-
  
---

# samsungcloudplatform_iam_group_policies (Data Source)



## Example Usage

```terraform
data "samsungcloudplatform_iam_groups" "my_own_groups" {
  group_name = "AdministratorGroup"
}

data "samsungcloudplatform_iam_group_policies" "my_group_policies" {
  group_id = data.samsungcloudplatform_iam_groups.my_own_groups.contents[0].group_id
}

output "result_my_groups" {
  value = data.samsungcloudplatform_iam_group_policies.my_group_policies
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `group_id` (String) Group ID

### Optional

- `filter` (Block Set) (see [below for nested schema](#nestedblock--filter))
- `policy_name` (String) Policy name
- `policy_type` (String) Policy type

### Read-Only

- `contents` (Block List) Contents list (see [below for nested schema](#nestedblock--contents))
- `id` (String) The ID of this resource.
- `total_count` (Number) Total count

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

- `created_by` (String) Creator's ID
- `created_by_email` (String) Creator's email
- `created_by_name` (String) Creator's name
- `created_dt` (String) Created date
- `description` (String) Description
- `modified_by` (String) Modifier's ID
- `modified_by_email` (String) Modifier's email
- `modified_by_name` (String) Modifier's name
- `modified_dt` (String) Modified date
- `policy_id` (String) Policy ID
- `policy_name` (String) Policy name
- `policy_type` (String) Policy type
- `principal_policy_id` (String) Principal policy ID


