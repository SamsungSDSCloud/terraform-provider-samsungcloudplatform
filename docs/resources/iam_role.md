---
page_title: "scp_iam_role Resource - scp"
subcategory: ""
description: |-
  
---

# Resource: scp_iam_role




## Example Usage

```terraform
resource "scp_iam_role" "my_role01" {
  role_name = var.name

  trust_principals {
    project_ids = [var.proj_id]
    user_srns = [var.srn]
  }

  tags {
    tag_key = "tk01"
    tag_value = "tv01"
  }
  tags {
    tag_key = "tk02"
    tag_value = "tv02"
  }

  description = ""
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `role_name` (String) Role name

### Optional

- `description` (String) Description
- `policy_ids` (Set of String) List of policy IDs
- `tags` (Block List) Tag list (see [below for nested schema](#nestedblock--tags))
- `trust_principals` (Block Set) Performing subjects (see [below for nested schema](#nestedblock--trust_principals))

### Read-Only

- `created_by` (String) Creator's ID
- `created_by_email` (String) Creator's email
- `created_by_name` (String) Creator's name
- `created_dt` (String) Created date
- `id` (String) The ID of this resource.
- `modified_by` (String) Modifier's ID
- `modified_by_email` (String) Modifier's email
- `modified_by_name` (String) Modifier's name
- `modified_dt` (String) Modified date
- `project_id` (String) Project ID
- `role_policy_count` (Number) Role's policy count
- `role_srn` (String) Role's SRN
- `session_time` (Number) Session time

<a id="nestedblock--tags"></a>
### Nested Schema for `tags`

Required:

- `tag_key` (String) Tag key

Optional:

- `tag_value` (String) Tag value


<a id="nestedblock--trust_principals"></a>
### Nested Schema for `trust_principals`

Optional:

- `project_ids` (List of String) Project IDs
- `user_srns` (List of String) User SRNs