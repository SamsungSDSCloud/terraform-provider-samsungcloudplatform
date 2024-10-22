data "samsungcloudplatform_bm_server" "server" {
  server_id = "BAREMETAL-XXXXXX"
}

output "output_bm_server" {
  value = data.samsungcloudplatform_bm_server.server
}
