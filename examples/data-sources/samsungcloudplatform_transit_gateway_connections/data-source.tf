data "samsungcloudplatform_transit_gateway_connections" "pjt_tgw_conns" {
}

output "contents" {
  value = data.samsungcloudplatform_transit_gateway_connections.pjt_tgw_conns.contents
}
