---
page_title: "samsungcloudplatform_iam_access_key Resource - samsungcloudplatform"
subcategory: ""
description: |-
  Provides IAM access key resource.
---

# Resource: samsungcloudplatform_iam_access_key

Provides IAM access key resource.


## Example Usage

```terraform
resource "samsungcloudplatform_iam_access_key" "my_access_key1" {
  project_id = var.project_id
  duration_days = var.duration_days
  access_key_activated = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `duration_days` (Number) Expiration time (days), set to zero to get permanent key
- `project_id` (String) Project ID

### Optional

- `access_key_activated` (Boolean) Access key activation

### Read-Only

- `access_key` (String) Access key
- `access_key_id` (String) Access key ID
- `access_key_state` (String) Access key state
- `access_secret_key` (String) Access secret key
- `expired_dt` (String) Expired date
- `id` (String) The ID of this resource.
- `project_name` (String) Project name


