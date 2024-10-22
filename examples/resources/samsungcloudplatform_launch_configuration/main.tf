resource "samsungcloudplatform_launch_configuration" "my_launch_configuration" {
  dynamic "block_storages" {
    for_each = var.block_storages
    content {
      block_storage_size = block_storages.value["block_storage_size"]
      disk_type = block_storages.value["disk_type"]
      encryption_enabled = block_storages.value["encryption_enabled"]
      is_boot_disk = block_storages.value["is_boot_disk"]
    }
  }
  image_id = var.image_id
  initial_script = var.initial_script
  key_pair_id = var.key_pair_id
  lc_name = var.lc_name
  server_type = var.server_type
  service_zone_id = var.service_zone_id
}
