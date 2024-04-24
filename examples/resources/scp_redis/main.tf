data "scp_region" "region" {
  filter {
    name = "location"
    values = ["KR-WEST-2"]
  }
}

data "scp_obs_storages" "obs_storage" {
  service_zone_id     = data.scp_region.region.id
  object_storage_name = "demo_object_storage_name"
}

data "scp_standard_image" "redis_7_2_0_image" {
  service_group = "DATABASE"
  service       = "Redis"
  region        = data.scp_region.region.location

  filter {
    name   = "image_name"
    values = ["Redis 7.2.0"]
  }
}

resource "scp_redis" "demo_db" {
  subnet_id = "SUBNET-123456789"
  security_group_ids = ["FIREWALL_SECURITY_GROUP-123456789", "FIREWALL_SECURITY_GROUP-987654321"]
  service_zone_id = data.scp_region.region.id

  image_id = data.scp_standard_image.redis_7_2_0_image.id
  contract_period = "1 Year"
  next_contract_period = "None"
  nat_enabled = true
  redis_name = "rediscluster"
  redis_state = "RUNNING"

  database_port = 6378
  database_user_password = ""

  encryption_enabled = true
  server_type = "redis1v1m2"
  timezone = "Asia/Seoul"

  redis_servers {
    redis_server_name = "demoredis-01"
    nat_public_ip_id = null
    server_role_type = "MASTER"
  }

  block_storages {
    block_storage_type = "SSD"
    block_storage_size = 50
  }

  backup  {
    object_storage_id = data.scp_obs_storages.obs_storage.contents[0].object_storage_id
    backup_retention_period = "15D"
    backup_start_hour = 7
  }
}
