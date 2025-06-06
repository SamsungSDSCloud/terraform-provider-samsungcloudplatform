---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "samsungcloudplatform_vpc_peering_detail Data Source - samsungcloudplatform"
subcategory: ""
description: |-
  Provides a VPC Peering detail.
---

# samsungcloudplatform_vpc_peering_detail (Data Source)

Provides a VPC Peering detail.

## Example Usage

```terraform
data "samsungcloudplatform_vpc_peerings" "peering" {
}

data "samsungcloudplatform_vpc_peering_detail" "detail" {
   vpc_peering_id = data.samsungcloudplatform_vpc_peerings.peering.contents[0].vpc_peering_id
}

output "detail" {
  value = data.samsungcloudplatform_vpc_peering_detail.detail
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `vpc_peering_id` (String) Vpc Peering Id

### Read-Only

- `approved_by` (String) Approved By
- `approved_dt` (String) Approved Date
- `approver_firewall_enabled` (Boolean) Approver Firewall Enabled
- `approver_project_id` (String) Approver Project Id
- `approver_vpc_id` (String) Approver Vpc Id
- `block_id` (String) Block Id
- `completed_dt` (String) Complated Date
- `created_by` (String) Created By
- `created_dt` (String) Created Date
- `id` (String) The ID of this resource.
- `modified_by` (String) Modified By
- `modified_dt` (String) Modified Date
- `product_group_id` (String) Product Group Id
- `project_id` (String) Project Id
- `requested_by` (String) Requested By
- `requested_dt` (String) Requested Date
- `requester_firewall_enabled` (Boolean) Requester Firewall Enabled
- `requester_project_id` (String) Requester Project Id
- `requester_vpc_id` (String) Requester Vpc Id
- `service_zone_id` (String) Service Zone Id
- `vpc_peering_description` (String) Vpc Peering Description
- `vpc_peering_name` (String) Vpc Peering Name
- `vpc_peering_state` (String) Vpc Peering State
- `vpc_peering_type` (String) Vpc Peering Type


