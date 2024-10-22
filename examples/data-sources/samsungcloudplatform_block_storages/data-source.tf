data "samsungcloudplatform_block_storages" "my_scp_block_storages" {
}

output "output_my_scp_block_storages" {
  value = data.samsungcloudplatform_block_storages.my_scp_block_storages
}
