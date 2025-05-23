---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "samsungcloudplatform_redis Data Source - samsungcloudplatform"
subcategory: ""
description: |-
  Search single redis database.
---

# samsungcloudplatform_redis (Data Source)

Search single redis database.

## Example Usage

```terraform
data "samsungcloudplatform_redis" "my_scp_redis" {
  redis_id = "SERVICE-123456789"
}

output "output_my_scp_redis" {
  value = data.samsungcloudplatform_redis.my_scp_redis
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `redis_id` (String) redis  Id

### Optional

- `security_group_ids` (List of String) security group ids

### Read-Only

- `backup_config` (Block List) backup config (see [below for nested schema](#nestedblock--backup_config))
- `block_id` (String) block id
- `contract` (Block List) contract (see [below for nested schema](#nestedblock--contract))
- `created_by` (String) created by
- `created_dt` (String) created dt
- `database_version` (String) database version
- `id` (String) The ID of this resource.
- `image_id` (String) image Id
- `maintenance` (Block List) maintenance (see [below for nested schema](#nestedblock--maintenance))
- `modified_by` (String) modified by
- `modified_dt` (String) modified dt
- `nat_ip_address` (String) nat ip address
- `project_id` (String) project id
- `redis_initial_config` (Block List) redis initial config (see [below for nested schema](#nestedblock--redis_initial_config))
- `redis_name` (String) redis  Name
- `redis_server_group` (Block List) redis server group (see [below for nested schema](#nestedblock--redis_server_group))
- `redis_state` (String) redis  State
- `sentinel_server` (Block List) redis server group (see [below for nested schema](#nestedblock--sentinel_server))
- `service_zone_id` (String) service zone id
- `subnet_id` (String) subnet Id
- `timezone` (String) timezone
- `vpc_id` (String) vPC Id

<a id="nestedblock--backup_config"></a>
### Nested Schema for `backup_config`

Read-Only:

- `full_backup_config` (Block List) full backup config (see [below for nested schema](#nestedblock--backup_config--full_backup_config))

<a id="nestedblock--backup_config--full_backup_config"></a>
### Nested Schema for `backup_config.full_backup_config`

Read-Only:

- `archive_backup_schedule_frequency` (String) archive backup schedule frequency
- `backup_retention_period` (String) backup retention period
- `backup_start_hour` (Number) backup start hour
- `object_storage_bucket_id` (String) object storage bucket id



<a id="nestedblock--contract"></a>
### Nested Schema for `contract`

Read-Only:

- `contract_end_date` (String) contract end date
- `contract_period` (String) contract period
- `contract_start_date` (String) contract start date
- `next_contract_end_date` (String) next contract end date
- `next_contract_period` (String) next contract period


<a id="nestedblock--maintenance"></a>
### Nested Schema for `maintenance`

Read-Only:

- `maintenance_period` (Number) maintenance period
- `maintenance_start_day_of_week` (String) maintenance start day of week
- `maintenance_start_time` (String) maintenance start time


<a id="nestedblock--redis_initial_config"></a>
### Nested Schema for `redis_initial_config`

Read-Only:

- `database_port` (Number) database port
- `sentinel_port` (Number) sentinel port


<a id="nestedblock--redis_server_group"></a>
### Nested Schema for `redis_server_group`

Read-Only:

- `block_storages` (Block List) block storages (see [below for nested schema](#nestedblock--redis_server_group--block_storages))
- `encryption_enabled` (Boolean) encryption enabled
- `redis_servers` (Block List) redis servers (see [below for nested schema](#nestedblock--redis_server_group--redis_servers))
- `server_group_role_type` (String) server group role type
- `server_type` (String) server type

<a id="nestedblock--redis_server_group--block_storages"></a>
### Nested Schema for `redis_server_group.block_storages`

Read-Only:

- `block_storage_group_id` (String) block storage group id
- `block_storage_name` (String) block storage name
- `block_storage_role_type` (String) block storage role type
- `block_storage_size` (Number) block Storage size
- `block_storage_type` (String) block storage type


<a id="nestedblock--redis_server_group--redis_servers"></a>
### Nested Schema for `redis_server_group.redis_servers`

Read-Only:

- `created_by` (String) created by
- `created_dt` (String) created dt
- `modified_by` (String) modified by
- `modified_dt` (String) modified dt
- `nat_public_ip_address` (String) nat public ip address
- `redis_server_id` (String) redis server id
- `redis_server_name` (String) redis server name
- `redis_server_state` (String) redis server state
- `server_role_type` (String) server role type
- `subnet_ip_address` (String) subnet ip address



<a id="nestedblock--sentinel_server"></a>
### Nested Schema for `sentinel_server`

Read-Only:

- `block_storages` (Block List) block storages (see [below for nested schema](#nestedblock--sentinel_server--block_storages))
- `created_by` (String) created by
- `created_dt` (String) created dt
- `encryption_enabled` (Boolean) encryption enabled
- `modified_by` (String) modified by
- `modified_dt` (String) modified dt
- `nat_public_ip_address` (String) nat public ip address
- `sentinel_server_id` (String) sentinel server id
- `sentinel_server_name` (String) sentinel server name
- `sentinel_server_state` (String) sentinel server state
- `server_type` (String) server type
- `subnet_ip_address` (String) subnet ip address

<a id="nestedblock--sentinel_server--block_storages"></a>
### Nested Schema for `sentinel_server.block_storages`

Read-Only:

- `block_storage_group_id` (String) block storage group id
- `block_storage_name` (String) block storage name
- `block_storage_role_type` (String) block storage role type
- `block_storage_size` (Number) block Storage size
- `block_storage_type` (String) block storage type


