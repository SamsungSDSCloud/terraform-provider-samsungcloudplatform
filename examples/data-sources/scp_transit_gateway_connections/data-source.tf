data "scp_transit_gateway_connections" "pjt_tgw_conns" {
}

output "contents" {
  value = data.scp_transit_gateway_connections.pjt_tgw_conns.contents
}
