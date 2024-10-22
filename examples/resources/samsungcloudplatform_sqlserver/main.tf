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

data "samsungcloudplatform_standard_image" "mssql_2019_enterprise_eng_image" {
  service_group = "DATABASE"
  service = "Microsoft SQL Server"
  region = data.samsungcloudplatform_region.region.location

  filter {
    name = "image_name"
    values = ["Microsoft SQL Server 2019 Enterprise ENG"]
  }
}

resource "samsungcloudplatform_sqlserver" "my_ms_sql" {
  subnet_id = "SUBNET-12345678"
  security_group_ids = ["FIREWALL_SECURITY_GROUP-12345678", "FIREWALL_SECURITY_GROUP-87654321"]
  service_zone_id = data.samsungcloudplatform_region.region.id

  sqlserver_servers = [
    {
      sqlserver_server_name = "sqlserver-pri"
      server_role_type = "PRIMARY"
    },
    {
      sqlserver_server_name = "sqlserver-sec"
      server_role_type = "SECONDARY"
    }
  ]

  image_id = data.samsungcloudplatform_standard_image.mssql_2019_enterprise_eng_image.id
  audit_enabled = true
  contract_period = "1 Year"
  next_contract_period = "None"
  nat_enabled = false
  nat_public_ip_id = null
  postgresql_cluster_name = "sqlservercluster"
  postgresql_cluster_state = "RUNNING"

  database_service_name = "MSsql"
  database_collation = "SQL_Latin1_General_CP1_CI_AS"
  license = "AAAAA-BBBBB-CCCCC-DDDDD-EEEEE"
  database_names = ["mssql_db","mssql_db2","mssql_db3"]
  database_port = 2866
  database_user_name = "dbuser"
  database_user_password = "pa$$w0rd"

  encryption_enabled = true
  server_type = "db1v2m4"
  timezone = "Asia/Seoul"

  block_storages = [
    {
      block_storage_type = "SSD"
      block_storage_size = 10
    }
  ]

  backup = [
    {
      object_storage_id = data.samsungcloudplatform_obs_storages.obs_storage.contents[0].object_storage_id
      archive_backup_schedule_frequency = "30M"
      backup_retention_period = "15D"
      backup_start_hour = 7
      full_backup_day_of_week = "MONDAY"
    }
  ]
}
