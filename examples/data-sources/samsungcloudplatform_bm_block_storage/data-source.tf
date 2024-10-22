data "samsungcloudplatform_bm_block_storage" "my_scp_bm_block_storage" {
  storage_id = "STORAGE-xxxxxxxxxxxxxxxxxxxxx"
}

output "output_my_scp_bm_block_storage_org"{
  value = data.samsungcloudplatform_bm_block_storage.my_scp_bm_block_storage
}
