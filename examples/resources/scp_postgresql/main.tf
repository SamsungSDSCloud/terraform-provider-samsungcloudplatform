data "scp_region" "region" {
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

resource "scp_postgresql" "my_pg_db" {
  image_id = data.scp_standard_image.postgres_13_6_image.id

  server_name_prefix = "pg-prefix"
  cluster_name       = "pgclusterxx"

  cpu_count          = var.cpu
  memory_size_gb     = var.memory

  contract_discount = "None"

  vpc_id             = data.terraform_remote_state.vpc.outputs.id
  subnet_id          = data.terraform_remote_state.subnet.outputs.id
  security_group_ids = [data.terraform_remote_state.security-group.outputs.id]

  db_name            = var.server_name
  db_user_id         = var.id
  db_user_password   = var.password
  db_port            = 2866

  timezone = "Asia/Seoul"

  data_disk_type = "SSD"
  data_storage_size_gb = 10

  additional_storage {
    product_name    = "SSD"
    storage_usage   = "DATA"
    storage_size_gb = 10
  }

  #high_availability {
  #  active_availability_zone_name  = "AZ1"
  #  standby_availability_zone_name = "AZ2"
  #}

  backup {
    backup_method = "s3api"
    retention_day = 7
    start_hour = 23
  }
}
