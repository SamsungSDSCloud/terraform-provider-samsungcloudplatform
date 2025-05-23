---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "samsungcloudplatform_iam_access_keys Data Source - samsungcloudplatform"
subcategory: ""
description: |-
  
---

# samsungcloudplatform_iam_access_keys (Data Source)



## Example Usage

```terraform
data "samsungcloudplatform_iam_access_keys" "my_access_keys" {

}

output "result_my_access_keys" {
  value = data.samsungcloudplatform_iam_access_keys.my_access_keys
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `access_key_project_type` (String) Access key's project type
- `access_key_state` (String) Access key state (ACTIVATED or DEACTIVATED)
- `active_yn` (Boolean) Whether the key is activated or not
- `filter` (Block Set) (see [below for nested schema](#nestedblock--filter))
- `project_id` (String) Project ID
- `project_name` (String) Access key's project name

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

- `access_key` (String) Access key
- `access_key_activated` (Boolean) Access key activated
- `access_key_id` (String) Access key ID
- `access_key_state` (String) Access key state
- `created_by` (String) Creator's ID
- `created_by_email` (String) Creator's email
- `created_by_name` (String) Creator's name
- `created_dt` (String) Created date
- `expired_dt` (String) Expiration date
- `modified_by` (String) Modifier's ID
- `modified_by_email` (String) Modifier's email
- `modified_by_name` (String) Modifier's name
- `modified_dt` (String) Modified date
- `project_id` (String) Project ID
- `project_name` (String) Project name
- `secret_vault_count` (Number) Secret Valut count


