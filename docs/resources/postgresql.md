---
page_title: "samsungcloudplatform_postgresql Resource - samsungcloudplatform"
subcategory: ""
description: |-
  Provides a PostgreSQL Database resource.
---

# Resource: samsungcloudplatform_postgresql

Provides a PostgreSQL Database resource.


## Example Usage

```terraform
data "samsungcloudplatform_region" "region" {
  filter {
    name = "location"
    values = ["KR-WEST-2"]
  }
}

data "samsungcloudplatform_obs_storages" "obs_storage" {
  service_zone_id     = data.samsungcloudplatform_region.region.id
  object_storage_name = "demo_object_storage_name"
}

data "samsungcloudplatform_standard_image" "postgres_13_6_image" {
  service_group = "DATABASE"
  service       = "PostgreSQL"
  region        = data.samsungcloudplatform_region.region.location

  filter {
    name   = "image_name"
    values = ["PostgreSQL Community 13.6"]
  }
}

resource "samsungcloudplatform_postgresql" "demo_db" {
  subnet_id = "SUBNET-123456789"
  security_group_ids = ["FIREWALL_SECURITY_GROUP-123456789", "FIREWALL_SECURITY_GROUP-987654321"]
  service_zone_id = data.samsungcloudplatform_region.region.id

  postgresql_servers {
    postgresql_server_name = "demopost-01"
    server_role_type = "ACTIVE"
  }

  image_id = data.samsungcloudplatform_standard_image.postgres_13_6_image.id
  audit_enabled = true
  contract_period = "1 Year"
  next_contract_period = "None"
  nat_enabled = true
  nat_public_ip_id = null
  postgresql_cluster_name = "demopostcluster"
  postgresql_cluster_state = "RUNNING"

  database_encoding = "UTF8"
  database_locale = "C"
  database_name = "demodb"
  database_port = 2866
  database_user_name = "demouser"
  database_user_password = ""

  encryption_enabled = true
  server_type = "db1v2m4"
  timezone = "Asia/Seoul"

  block_storages {
    block_storage_type = "SSD"
    block_storage_role_type = "DATA"
    block_storage_size = 10
  }

  backup {
    object_storage_id = data.samsungcloudplatform_obs_storages.obs_storage.contents[0].object_storage_id
    archive_backup_schedule_frequency = "30M"
    backup_retention_period = "15D"
    backup_start_hour = 7
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `audit_enabled` (Boolean) Whether to use database audit logging.
- `block_storages` (Block List, Min: 1, Max: 10) block storage. (see [below for nested schema](#nestedblock--block_storages))
- `contract_period` (String) Contract (None|1 Year|3 Year)
- `database_encoding` (String) Postgresql encoding. (Only 'UTF8' for now)
- `database_locale` (String) Postgresql locale. (Only 'C' for now)
- `database_name` (String) Name of database. (only English alphabets or numbers between 3 and 20 characters)
- `database_port` (Number) Port number of database. (1024 to 65535)
- `database_user_name` (String) User account id of database. (2 to 20 lowercase alphabets)
- `database_user_password` (String, Sensitive) User account password of database.
- `encryption_enabled` (Boolean) Whether to use storage encryption.
- `image_id` (String) Postgresql virtual server image id.
- `postgresql_cluster_name` (String) Name of database cluster. (3 to 20 characters only)
- `postgresql_cluster_state` (String) postgresql cluster state (RUNNING|STOPPED)
- `postgresql_servers` (Block List, Min: 1, Max: 2) postgresql servers (HA configuration when entering two server specifications) (see [below for nested schema](#nestedblock--postgresql_servers))
- `security_group_ids` (List of String) Security-Group ids of this postgresql DB. Each security-group must be a valid security-group resource which is attached to the VPC.
- `server_type` (String) Server type
- `service_zone_id` (String) Service Zone Id
- `subnet_id` (String) Subnet id of this database server. Subnet must be a valid subnet resource which is attached to the VPC.
- `timezone` (String) Timezone setting of this database.

### Optional

- `backup` (Block Set, Max: 1) (see [below for nested schema](#nestedblock--backup))
- `nat_enabled` (Boolean) Whether to use nat.
- `nat_public_ip_id` (String) Public IP for NAT. If it is null, it is automatically allocated.
- `next_contract_period` (String) Next contract (None|1 Year|3 Year)
- `tags` (Map of String)
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `id` (String) The ID of this resource.
- `nat_ip_address` (String) nat ip address
- `virtual_ip_address` (String) virtual ip address
- `vpc_id` (String) vpc id

<a id="nestedblock--block_storages"></a>
### Nested Schema for `block_storages`

Required:

- `block_storage_role_type` (String) Storage usage. (DATA|ARCHIVE|TEMP|BACKUP)
- `block_storage_size` (Number) Block Storage Size (10 to 5120)
- `block_storage_type` (String) Storage product name. (SSD|HDD)

Read-Only:

- `block_storage_group_id` (String) Block storage group id


<a id="nestedblock--postgresql_servers"></a>
### Nested Schema for `postgresql_servers`

Required:

- `postgresql_server_name` (String) Postgresql database server names. (3 to 20 lowercase and number with dash and the first character should be an lowercase letter.)
- `server_role_type` (String) Server role type Enter 'ACTIVE' for a single server configuration. (ACTIVE | STANDBY)

Optional:

- `availability_zone_name` (String) Availability Zone Name. The single server does not input anything. (AZ1|AZ2)


<a id="nestedblock--backup"></a>
### Nested Schema for `backup`

Required:

- `archive_backup_schedule_frequency` (String) Backup File Schedule Frequency.(5M|10M|30M|1H)
- `backup_retention_period` (String) Backup File Retention Day.(7D <= day <= 35D)
- `backup_start_hour` (Number) The time at which the backup starts. (from 0 to 23)

Optional:

- `object_storage_id` (String) Object storage ID where backup files will be stored.


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `update` (String)


