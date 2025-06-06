---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "samsungcloudplatform_mysql Data Source - samsungcloudplatform"
subcategory: ""
description: |-
  Search single mysql database.
---

# samsungcloudplatform_mysql (Data Source)

Search single mysql database.

## Example Usage

```terraform
data "samsungcloudplatform_mysql" "my_scp_mysql" {
  mysql_cluster_id = "SERVICE-123456789"
}

output "output_my_scp_mysql" {
  value = data.samsungcloudplatform_mysql.my_scp_mysql
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `mysql_cluster_id` (String) mysql Cluster Id

### Optional

- `mysql_replica_cluster_ids` (List of String) mysql replica cluster ids
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
- `mysql_cluster_name` (String) mysql Cluster Name
- `mysql_cluster_state` (String) mysql Cluster State
- `mysql_initial_config` (Block List) mysql initial config (see [below for nested schema](#nestedblock--mysql_initial_config))
- `mysql_master_cluster_id` (String) mysql master cluster id
- `mysql_server_group` (Block List) mysql server group (see [below for nested schema](#nestedblock--mysql_server_group))
- `nat_ip_address` (String) nat ip address
- `project_id` (String) project id
- `service_zone_id` (String) service zone id
- `subnet_id` (String) subnet Id
- `timezone` (String) timezone
- `vpc_id` (String) vPC Id

<a id="nestedblock--backup_config"></a>
### Nested Schema for `backup_config`

Read-Only:

- `full_backup_config` (Block List) full backup config (see [below for nested schema](#nestedblock--backup_config--full_backup_config))
- `incremental_backup_config` (Block List) incremental_backup_config (see [below for nested schema](#nestedblock--backup_config--incremental_backup_config))

<a id="nestedblock--backup_config--full_backup_config"></a>
### Nested Schema for `backup_config.full_backup_config`

Read-Only:

- `archive_backup_schedule_frequency` (String) archive backup schedule frequency
- `backup_retention_period` (String) backup retention period
- `backup_start_hour` (Number) backup start hour
- `object_storage_bucket_id` (String) object storage bucket id


<a id="nestedblock--backup_config--incremental_backup_config"></a>
### Nested Schema for `backup_config.incremental_backup_config`

Read-Only:

- `archive_backup_schedule_frequency` (String) archive backup schedule frequency
- `backup_retention_period` (String) backup retention period
- `backup_schedule_frequency` (String) backup schedule frequency
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


<a id="nestedblock--mysql_initial_config"></a>
### Nested Schema for `mysql_initial_config`

Read-Only:

- `database_character_set` (String) database character set
- `database_name` (String) database name
- `database_port` (Number) database port
- `database_user_name` (String) database user name


<a id="nestedblock--mysql_server_group"></a>
### Nested Schema for `mysql_server_group`

Read-Only:

- `block_storages` (Block List) block storages (see [below for nested schema](#nestedblock--mysql_server_group--block_storages))
- `encryption_enabled` (Boolean) encryption enabled
- `mysql_servers` (Block List) mysql servers (see [below for nested schema](#nestedblock--mysql_server_group--mysql_servers))
- `server_group_role_type` (String) server group role type
- `server_type` (String) server type
- `virtual_ip_address` (String) virtual ip address

<a id="nestedblock--mysql_server_group--block_storages"></a>
### Nested Schema for `mysql_server_group.block_storages`

Read-Only:

- `block_storage_group_id` (String) block storage group id
- `block_storage_name` (String) block storage name
- `block_storage_role_type` (String) block storage role type
- `block_storage_size` (Number) block Storage size
- `block_storage_type` (String) block storage type


<a id="nestedblock--mysql_server_group--mysql_servers"></a>
### Nested Schema for `mysql_server_group.mysql_servers`

Read-Only:

- `availability_zone_name` (String) availability zone name
- `created_by` (String) created by
- `created_dt` (String) created dt
- `modified_by` (String) modified by
- `modified_dt` (String) modified dt
- `mysql_server_id` (String) mysql server id
- `mysql_server_name` (String) mysql server name
- `mysql_server_state` (String) mysql server state
- `server_role_type` (String) server role type
- `subnet_ip_address` (String) subnet ip address


