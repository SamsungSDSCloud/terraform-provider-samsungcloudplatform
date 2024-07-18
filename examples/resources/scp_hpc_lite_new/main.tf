resource "scp_hpc_lite_new" "hpc_lite_new" {
  co_service_zone_id = var.co_service_zone_id
  contract = var.contract
  hyper_threading_enabled = var.hyper_threading_enabled
  image_id = var.image_id
  init_script = var.init_script
  os_user_id = var.os_user_id
  os_user_password = var.os_user_password
  product_group_id = var.product_group_id
  resource_pool_id = var.resource_pool_id
  server_type = var.server_type
  service_zone_id = var.service_zone_id
  vlan_pool_cidr = var.vlan_pool_cidr

  dynamic "server_details" {
    for_each = var.server_details
    content {
      server_name = server_details.value.server_name
      ip_address = try(server_details.value.ip_address, null)
    }
  }
  tags = {
    tk01 = "tv01"
  }
}
