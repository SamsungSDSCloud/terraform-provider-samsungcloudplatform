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

data "samsungcloudplatform_standard_image" "mysql_image" {
  service_group = "DATABASE"
  service       = "MySQL"
  region        = data.samsungcloudplatform_region.region.location

  filter {
    name   = "image_name"
    values = ["MySQL 8.0.32"]
  }
}

resource "samsungcloudplatform_mysql" "demo_db" {
  subnet_id = "SUBNET-123456789"
  security_group_ids = ["FIREWALL_SECURITY_GROUP-123456789", "FIREWALL_SECURITY_GROUP-987654321"]
  service_zone_id = data.samsungcloudplatform_region.region.id

  mysql_servers {
    mysql_server_name = "demomysql-01"
    server_role_type = "ACTIVE"
  }

  image_id = data.samsungcloudplatform_standard_image.mysql_image.id
  contract_period = "1 Year"
  next_contract_period = "None"
  nat_enabled = true
  nat_public_ip_id = null
  mysql_cluster_name = "demomysql"
  mysql_cluster_state = "RUNNING"

  database_case_sensitivity = false
  database_character_set = "utf8mb3"
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
    object_storage_id = data.samsungcloudplatform_obs_storages.obs_storage.contents[0].object_storage_id
    archive_backup_schedule_frequency = "30M"
    backup_retention_period = "15D"
    backup_start_hour = 7
  }
}
