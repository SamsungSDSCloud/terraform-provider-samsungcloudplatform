---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "samsungcloudplatform_auto_scaling_group_policies Data Source - scp"
subcategory: ""
description: |-
  Provides list of Auto-Scaling Group policies
---

# samsungcloudplatform_auto_scaling_group_policies (Data Source)

Provides list of Auto-Scaling Group policies

## Example Usage

```terraform
# Find all Auto-Scaling Group policies
data "samsungcloudplatform_auto_scaling_group_policies" "my_auto_scaling_group_policies1" {
  asg_id = "AUTO_SCALING_GROUP-XXXXX"
}

# Find all Auto-Scaling Group policies
data "samsungcloudplatform_auto_scaling_group_policies" "my_auto_scaling_group_policies2" {
  asg_id = "AUTO_SCALING_GROUP-XXXXX"

  # Sort in ascending order of creation date
  sort = "createdDt:asc"

  # Apply filter for 'policy_name' regex value "test"
  filter {
    name = "policy_name"
    values = ["test"]
    use_regex = true
  }
}

output "output_scp_auto_scaling_group_policies1" {
  value = data.samsungcloudplatform_auto_scaling_group_policies.my_auto_scaling_group_policies1
}

output "output_scp_auto_scaling_group_policies2" {
  value = data.samsungcloudplatform_auto_scaling_group_policies.my_auto_scaling_group_policies2
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `asg_id` (String) Auto-Scaling Group ID

### Optional

- `filter` (Block Set) (see [below for nested schema](#nestedblock--filter))
- `metric_method` (String) Metric method
- `metric_type` (String) Metric type
- `page` (Number) Page start number from which to get the list
- `policy_name` (String) Policy name
- `scale_type` (String) Scale type
- `size` (Number) Size to get list
- `sort` (String) Sort

### Read-Only

- `contents` (List of Object) Auto-Scaling Group policy list (see [below for nested schema](#nestedatt--contents))
- `id` (String) The ID of this resource.
- `total_count` (Number) Total list size

<a id="nestedblock--filter"></a>
### Nested Schema for `filter`

Required:

- `name` (String) Filtering target name
- `values` (List of String) Filtering values. Each matching value is appended. (OR rule)

Optional:

- `use_regex` (Boolean) Enable regex match for values


<a id="nestedatt--contents"></a>
### Nested Schema for `contents`

Read-Only:

- `asg_id` (String)
- `comparison_operator` (String)
- `cooldown_seconds` (Number)
- `created_by` (String)
- `created_dt` (String)
- `evaluation_minutes` (Number)
- `metric_method` (String)
- `metric_type` (String)
- `modified_by` (String)
- `modified_dt` (String)
- `policy_id` (String)
- `policy_name` (String)
- `policy_state` (String)
- `scale_method` (String)
- `scale_type` (String)
- `scale_value` (Number)
- `threshold` (String)
- `threshold_unit` (String)

