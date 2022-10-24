data "scp_file_storages" "my_scp_file_storages" {
}

output "output_my_scp_file_storages" {
  value = data.scp_file_storages.my_scp_file_storages
}
