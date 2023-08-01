data "scp_bm_servers" "servers" {
}

output "output_my_scp_block_storages" {
  value = data.scp_bm_servers.servers
}
