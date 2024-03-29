---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "scp_obs_bucket Data Source - scp"
subcategory: ""
description: |-
  Provides Object Bucket Info.
---

# scp_obs_bucket (Data Source)

Provides Object Bucket Info.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `obs_bucket_id` (String) Obs Bucket Id

### Optional

- `obs_urls` (Map of String) Obs Urls

### Read-Only

- `created_by` (String) Created By
- `created_dt` (String) Created Date
- `id` (String) The ID of this resource.
- `is_obs_bucket_dr_enabled` (Boolean) Dr Enabled
- `is_obs_bucket_ip_address_filter_enabled` (Boolean) Ip Filter Enabled
- `is_obs_object_creation_enabled` (Boolean) Object Creation Enabled
- `is_obs_system_bucket_enabled` (Boolean) System Bucket Enabled
- `is_replication_in_progress` (Boolean) Replication In Progress
- `modified_by` (String) Modified By
- `modified_dt` (String) Modified Date
- `multi_az_yn` (String) Multi Az Y/N
- `obs_bucket_access_ip_address_ranges` (List of Object) Bucket Access Ip Ranges (see [below for nested schema](#nestedatt--obs_bucket_access_ip_address_ranges))
- `obs_bucket_access_url` (String) Bucket Access Url
- `obs_bucket_dr_type` (String) Dr Type
- `obs_bucket_file_encryption_algorithm` (String) Bucket Encryption Algorithm
- `obs_bucket_file_encryption_enabled` (Boolean) Is Encryption Enabled
- `obs_bucket_file_encryption_type` (String) Bucket Encryption Type
- `obs_bucket_name` (String) Bucket Name
- `obs_bucket_state` (String) Bucket State
- `obs_bucket_used_size` (Number) Bucket Used Size
- `obs_bucket_used_type` (String) Bucket Used Type
- `obs_bucket_version_enabled` (Boolean) Versioning Enabled
- `obs_id` (String) Object Storage Id
- `obs_name` (String) Object Storage Name
- `obs_quota_id` (String) Obs Quota Id
- `obs_sync_bucket_id` (String) Obs Quota Name
- `obs_sync_bucket_name` (String) Obs Tenant Name
- `obs_sync_bucket_obs_name` (String) Pool Region
- `obs_sync_bucket_region` (String) System Id
- `obs_sync_bucket_zone_name` (String) Obs Quota Name
- `project_id` (String) Project Id
- `project_name` (String) Project Name
- `region` (String) Region
- `system_id` (String) System Id
- `system_name` (String) System Name
- `zone_id` (String) Zone Id
- `zone_name` (String) Zone Name

<a id="nestedatt--obs_bucket_access_ip_address_ranges"></a>
### Nested Schema for `obs_bucket_access_ip_address_ranges`

Read-Only:

- `ip_address_range` (String)
- `type` (String)


