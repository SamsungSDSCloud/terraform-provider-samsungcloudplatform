data "samsungcloudplatform_transit_gateways" "pjt_tgws" {
}

output "contents" {
  value = data.samsungcloudplatform_transit_gateways.pjt_tgws.contents
}
