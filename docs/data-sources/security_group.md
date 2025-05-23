---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "samsungcloudplatform_security_group Data Source - samsungcloudplatform"
subcategory: ""
description: |-
  Provides Security Group Info
---

# samsungcloudplatform_security_group (Data Source)

Provides Security Group Info

## Example Usage

```terraform
data "samsungcloudplatform_security_group" "my_sg" {
  security_group_id = "FIREWALL_SECURITY_GROUP-XXXXXXXXXXXXXXXXXXXXXX"
}

output "output_my_scp_sg" {
  value = data.samsungcloudplatform_security_group.my_sg
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `security_group_id` (String) Security Group ID

### Read-Only

- `created_by` (String) creator
- `created_dt` (String) created datetime
- `id` (String) The ID of this resource.
- `is_loggable` (Boolean) Is loggable
- `modified_by` (String) last modified user
- `modified_dt` (String) Resource modified datetime
- `project_id` (String) Project ID
- `rule_count` (Number) The number of Rules
- `scope` (String) Security Group Scope of Use
- `security_group_description` (String) Security Group description
- `security_group_name` (String) Security Group name
- `security_group_state` (String) Security Group state
- `vendor_object_id` (String) Vendor Object ID
- `vpc_id` (String) VPC ID
- `zone_id` (String) Service Zone ID


