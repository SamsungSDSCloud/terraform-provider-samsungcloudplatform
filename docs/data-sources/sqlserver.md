---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "samsungcloudplatform_sqlserver Data Source - samsungcloudplatform"
subcategory: ""
description: |-
  Search Detail MS SQL Server database.
---

# samsungcloudplatform_sqlserver (Data Source)

Search Detail MS SQL Server database.

## Example Usage

```terraform
data "samsungcloudplatform_sqlserver" "my_scp_sqlserver" {
  sqlserver_cluster_id = "SERVICE-123456789"
}

output "output_my_scp_sqlserver" {
  value = data.samsungcloudplatform_sqlserver.my_scp_sqlserver
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `sqlserver_cluster_id` (String) MS SQL Server Cluster Id

### Optional

- `security_group_ids` (List of String) Security group ids
- `sqlserver_initial_config` (Block List) MS SQL Server initial config (see [below for nested schema](#nestedblock--sqlserver_initial_config))

### Read-Only

- `audit_enabled` (Boolean) audit enabled
- `backup_config` (Block List) backup config (see [below for nested schema](#nestedblock--backup_config))
- `block_id` (String) Block id
- `contract` (Block List) contract (see [below for nested schema](#nestedblock--contract))
- `created_by` (String) created by
- `created_dt` (String) created dt
- `database_version` (String) MS SQL Server version
- `id` (String) The ID of this resource.
- `image_id` (String) Image Id
- `maintenance` (Block List) maintenance (see [below for nested schema](#nestedblock--maintenance))
- `modified_by` (String) modified by
- `modified_dt` (String) modified dt
- `nat_ip_address` (String) nat ip address
- `project_id` (String) project id
- `quorum_server_group` (Block List) MS SQL Server quorum server group (see [below for nested schema](#nestedblock--quorum_server_group))
- `service_zone_id` (String) service zone id
- `sqlserver_cluster_name` (String) MS SQL Server Cluster Name
- `sqlserver_cluster_state` (String) MS SQL Server Cluster State
- `sqlserver_master_cluster_id` (String) MS SQL Server master cluster id
- `sqlserver_secondary_cluster_id` (String) MS SQL Server secondary cluster id
- `sqlserver_server_group` (Block List) MS SQL Server server group (see [below for nested schema](#nestedblock--sqlserver_server_group))
- `subnet_id` (String) Subnet Id
- `timezone` (String) Timezone
- `vpc_id` (String) VPC Id

<a id="nestedblock--sqlserver_initial_config"></a>
### Nested Schema for `sqlserver_initial_config`

Required:

- `database_service_name` (String) MS SQL Server Database Service name

Optional:

- `database_names` (List of String) Database Name List
- `sqlserver_active_directory` (Block Set) MS SQL Server Active directory (see [below for nested schema](#nestedblock--sqlserver_initial_config--sqlserver_active_directory))

Read-Only:

- `database_collation` (String) Commands that specify how to sort and compare data
- `database_port` (Number) database port
- `database_user_name` (String) database user name

<a id="nestedblock--sqlserver_initial_config--sqlserver_active_directory"></a>
### Nested Schema for `sqlserver_initial_config.sqlserver_active_directory`

Optional:

- `dns_server_ips` (List of String) Active Directory DNS Server IPs

Read-Only:

- `ad_server_user_id` (String) Active Directory Server User ID
- `ad_server_user_password` (String, Sensitive) Active Directory Server User password
- `domain_name` (String) Active Directory Domain name
- `domain_net_bios_name` (String) Active Directory NetBios name
- `failover_cluster_name` (String) Active Directory Failover Cluster name



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
- `full_backup_day_of_week` (String) Full backup schedule(Day).
- `object_storage_bucket_id` (String) object storage bucket id



<a id="nestedblock--contract"></a>
### Nested Schema for `contract`

Read-Only:

- `contract_end_date` (String) Contract end date
- `contract_period` (String) Contract period
- `contract_start_date` (String) Contract start date
- `next_contract_end_date` (String) Next contract end date
- `next_contract_period` (String) Next contract period


<a id="nestedblock--maintenance"></a>
### Nested Schema for `maintenance`

Read-Only:

- `maintenance_period` (Number) maintenance period
- `maintenance_start_day_of_week` (String) maintenance start day of week
- `maintenance_start_time` (String) maintenance start time


<a id="nestedblock--quorum_server_group"></a>
### Nested Schema for `quorum_server_group`

Read-Only:

- `block_storages` (Block List) block storages (see [below for nested schema](#nestedblock--quorum_server_group--block_storages))
- `encryption_enabled` (Boolean) encryption enabled
- `server_group_role_type` (String) server group role type
- `server_type` (String) server type
- `sqlserver_servers` (Block List) MS SQL Server quorum servers (see [below for nested schema](#nestedblock--quorum_server_group--sqlserver_servers))
- `virtual_ip_address` (String) virtual ip address

<a id="nestedblock--quorum_server_group--block_storages"></a>
### Nested Schema for `quorum_server_group.block_storages`

Read-Only:

- `block_storage_group_id` (String) block storage group id
- `block_storage_name` (String) block storage name
- `block_storage_role_type` (String) block storage role type
- `block_storage_size` (Number) block Storage size
- `block_storage_type` (String) block storage type


<a id="nestedblock--quorum_server_group--sqlserver_servers"></a>
### Nested Schema for `quorum_server_group.sqlserver_servers`

Read-Only:

- `availability_zone_name` (String) availability zone name
- `created_by` (String) created by
- `created_dt` (String) created dt
- `modified_by` (String) modified by
- `modified_dt` (String) modified dt
- `server_role_type` (String) server role type
- `sqlserver_server_id` (String) MS SQL Server quorum server id
- `sqlserver_server_name` (String) MS SQL Server quorum server name
- `sqlserver_server_state` (String) MS SQL Server quorum server state
- `subnet_ip_address` (String) subnet ip address



<a id="nestedblock--sqlserver_server_group"></a>
### Nested Schema for `sqlserver_server_group`

Read-Only:

- `block_storages` (Block List) block storages (see [below for nested schema](#nestedblock--sqlserver_server_group--block_storages))
- `encryption_enabled` (Boolean) encryption enabled
- `server_group_role_type` (String) server group role type
- `server_type` (String) server type
- `sqlserver_servers` (Block List) MS SQL Server servers (see [below for nested schema](#nestedblock--sqlserver_server_group--sqlserver_servers))
- `virtual_ip_address` (String) virtual ip address

<a id="nestedblock--sqlserver_server_group--block_storages"></a>
### Nested Schema for `sqlserver_server_group.block_storages`

Read-Only:

- `block_storage_group_id` (String) block storage group id
- `block_storage_name` (String) block storage name
- `block_storage_role_type` (String) block storage role type
- `block_storage_size` (Number) block Storage size
- `block_storage_type` (String) block storage type


<a id="nestedblock--sqlserver_server_group--sqlserver_servers"></a>
### Nested Schema for `sqlserver_server_group.sqlserver_servers`

Read-Only:

- `availability_zone_name` (String) availability zone name
- `created_by` (String) created by
- `created_dt` (String) created dt
- `modified_by` (String) modified by
- `modified_dt` (String) modified dt
- `server_role_type` (String) server role type
- `sqlserver_server_id` (String) MS SQL Server server id
- `sqlserver_server_name` (String) MS SQL Server server name
- `sqlserver_server_state` (String) MS SQL Server server state
- `subnet_ip_address` (String) subnet ip address


