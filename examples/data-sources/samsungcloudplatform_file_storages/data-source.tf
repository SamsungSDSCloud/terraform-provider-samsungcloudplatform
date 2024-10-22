data "samsungcloudplatform_file_storages" "my_scp_file_storages" {
  file_storage_states = [
    "ACTIVE",
    "ERROR"
  ]
  sort = [
    "fileStorageName:DESC"
  ]
}

output "output_my_scp_file_storages" {
  value = data.samsungcloudplatform_file_storages.my_scp_file_storages
}
