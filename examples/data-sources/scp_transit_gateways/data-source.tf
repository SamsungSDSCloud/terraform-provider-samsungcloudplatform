data "scp_transit_gateways" "pjt_tgws" {
}

output "contents" {
  value = data.scp_transit_gateways.pjt_tgws.contents
}
