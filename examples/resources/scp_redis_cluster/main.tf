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

data "scp_standard_image" "redis_cluster_7_2_image" {
  service_group = "DATABASE"
  service       = "Redis"
  region        = data.scp_region.region.location

  filter {
    name   = "image_name"
    values = ["Redis Cluster 7.2.0"]
  }
}

resource "scp_redis_cluster" "demo_db" {
  subnet_id = "SUBNET-12345678"
  security_group_ids = ["FIREWALL_SECURITY_GROUP-12345678", "FIREWALL_SECURITY_GROUP-87654321"]
  service_zone_id = "ZONE-12345678"

  redis_servers = [
    {
      redis_server_name = "terraabc-01"
      nat_public_ip_id = null
      server_role_type = "MASTER"
    },
    {
      redis_server_name = "terraabc-02"
      nat_public_ip_id = null
      server_role_type = "MASTER"
    },
    {
      redis_server_name = "terraabc-03"
      nat_public_ip_id = null
      server_role_type = "MASTER"
    },
    {
      redis_server_name = "terraabc-04"
      nat_public_ip_id = null
      server_role_type = "REPLICA"
    },
    {
      redis_server_name = "terraabc-05"
      nat_public_ip_id = null
      server_role_type = "REPLICA"
    },
    {
      redis_server_name = "terraabc-06"
      nat_public_ip_id = null
      server_role_type = "REPLICA"
    }
  ]


  image_id = data.scp_standard_image.redis_cluster_7_2_image.id
  contract_period = "1 Year"
  next_contract_period = "None"
  nat_enabled = false
  redis_cluster_name = "democluster"
  redis_cluster_state = "RUNNING"

  database_port = 6378
  database_user_password = ""

  shards_count = 3
  shards_replica_count = 1

  encryption_enabled = true
  server_type = "redis1v1m2"
  timezone = "Asia/Seoul"

  block_storages = [
    {
      block_storage_role_type = "DATA"
      block_storage_type = "SSD"
      block_storage_size = 50
    }
  ]

  backup = [
    {
      object_storage_id = "S3OBJECTSTORAGE-12345678"
      backup_retention_period = "15D"
      backup_start_hour = 7
    }
  ]
}
