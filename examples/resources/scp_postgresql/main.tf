data "scp_region" "region" {
  filter {
    name = "location"
    values = ["KR-EAST-1"]
  }
}

data "scp_obs_storages" "obs_storage" {
  service_zone_id     = data.scp_region.region.id
  object_storage_name = "demo_object_storage_name"
}

data "scp_standard_image" "postgres_13_6_image" {
  service_group = "DATABASE"
  service       = "PostgreSQL"
  region        = data.scp_region.region.location

  filter {
    name   = "image_name"
    values = ["PostgreSQL Community 13.6"]
  }
}

resource "scp_postgresql" "demo_db" {
  subnet_id = "SUBNET-123456789"
  security_group_ids = ["FIREWALL_SECURITY_GROUP-123456789", "FIREWALL_SECURITY_GROUP-987654321"]
  service_zone_id = data.scp_region.region.id

  postgresql_servers {
    postgresql_server_name = "demoserver-01"
    server_role_type = "ACTIVE"
  }

  image_id = data.scp_standard_image.postgres_13_6_image.id
  audit_enabled = true
  contract_period = "1 Year"
  next_contract_period = "None"
  nat_enabled = false
  nat_public_ip_id = null
  postgresql_cluster_name = "democluster"
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

  backup  {
    object_storage_id = data.scp_obs_storages.obs_storage.contents[0].object_storage_id
    archive_backup_schedule_frequency = "30M"
    backup_retention_period = "15D"
    backup_start_hour = 7
  }
}
