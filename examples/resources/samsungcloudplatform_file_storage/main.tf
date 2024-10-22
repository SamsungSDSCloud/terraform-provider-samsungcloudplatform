resource "samsungcloudplatform_file_storage" "my_scp_file_storage" {
  file_storage_name = "fs_cifs_test"
  disk_type         = "HDD"
  file_storage_protocol = "CIFS"
  cifs_password = var.password
  product_names = [
    "HDD"
  ]
  service_zone_id = "ZONE-XXXXX"
}
