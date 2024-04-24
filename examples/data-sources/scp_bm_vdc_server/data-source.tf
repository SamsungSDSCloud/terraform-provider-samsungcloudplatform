data "scp_bm_vdc_server" "server" {
  server_id = "BAREMETALVDC-XXXXXX"
}

output "output_bm_server" {
  value = data.scp_bm_vdc_server.server
}
