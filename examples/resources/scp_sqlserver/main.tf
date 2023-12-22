data "scp_region" "region" {
}

data "scp_standard_image" "mssql_2019_enterprise_eng_image" {
  service_group = "DATABASE"
  service = "Microsoft SQL Server"
  region = data.scp_region.region.location

  filter {
    name = "image_name"
    values = ["Microsoft SQL Server 2019 Enterprise ENG"]
  }
}

resource "scp_sqlserver" "my_ms_sql" {
  image_id = data.scp_standard_image.mssql_2019_enterprise_eng_image.id

  server_group_name = "mssqlsvgr"
  virtual_server_name_prefix = "mssqlvs"

  vpc_id = data.terraform_remote_state.vpc.outputs.id
  subnet_id = data.terraform_remote_state.subnet.outputs.id
  security_group_ids = [data.terraform_remote_state.security-group.outputs.id]

  db_service_name = "Dbsvcname"
  db_name = "dba"
  db_user_id = var.id
  db_user_password = var.password
  db_port = 9548

  license_key = var.license

  cpu_count = var.cpu
  memory_size_gb = var.memory

  contract_discount = "None"

  timezone = "Asia/Seoul"

  db_collation = "Korean_Wansung_CS_AS"

  data_block_storage_size_gb = 100
  encrypt_enabled = false

  additional_block_storages {
    storage_usage = "DATA"
    storage_size_gb = 10
  }

  additional_db = ["dbb"]

  backup {
    backup_retention_day = 7
    backup_start_hour = 23
  }
}
