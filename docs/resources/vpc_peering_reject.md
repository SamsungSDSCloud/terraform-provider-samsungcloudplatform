---
page_title: "samsungcloudplatform_vpc_peering_reject Resource - samsungcloudplatform"
subcategory: ""
description: |-
  Reject Peering Request.
---

# Resource: samsungcloudplatform_vpc_peering_reject

Reject Peering Request.


## Example Usage

```terraform
data "samsungcloudplatform_vpc_peerings" "peerings" {
}

resource "samsungcloudplatform_vpc_peering_reject" "reject" {
  vpc_peering_id = "VPC_PEERING-XXXX"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `vpc_peering_id` (String) Vpc Peering Id

### Read-Only

- `id` (String) The ID of this resource.
- `vpc_peering_state` (String) Vpc Peering Id


