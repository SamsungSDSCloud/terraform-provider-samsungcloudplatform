data "scp_bm_block_storages" "my_scp_bm_block_storages" {
}

output "output_my_scp_bm_block_storages_org" {
  value = data.scp_bm_block_storages.my_scp_bm_block_storages
}
