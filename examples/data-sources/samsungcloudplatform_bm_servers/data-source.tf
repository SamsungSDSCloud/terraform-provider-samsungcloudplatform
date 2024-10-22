data "samsungcloudplatform_bm_servers" "servers" {
}

output "output_my_scp_block_storages" {
  value = data.samsungcloudplatform_bm_servers.servers
}
