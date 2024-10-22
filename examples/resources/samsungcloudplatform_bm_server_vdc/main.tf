data "samsungcloudplatform_project" project {}

resource "samsungcloudplatform_bm_server_vdc" "vdc_server" {

  # use default block_id and service_zone_id
  block_id = data.samsungcloudplatform_project.project.service_zones[0].block_id
  service_zone_id = data.samsungcloudplatform_project.project.service_zones[0].service_zone_id

  vdc_id = var.vdc_id
  subnet_id = var.subnet_id
  image_name = var.image_name

  contract_discount = var.contract_discount
  delete_protection = var.delete_protection
  admin_account = var.admin
  admin_password = var.password
  initial_script = var.initial_script

  cpu_count = var.cpu
  memory_size = var.memory

  dynamic "servers" {
    for_each = local.servers
    content {
      name = servers.value.name
      ip_address = servers.value.ip_address
      use_hyper_threading = servers.value.use_hyper_threading
      dns_enabled = servers.value.dns_enabled
    }
  }

  block_storages {
    name = "test-bs"
    product_name = "SSD"
    storage_size_gb = 10
    encrypted = false
  }

  timeouts {
    create = "40m"
    update = "40m"
    delete = "40m"
  }
}
