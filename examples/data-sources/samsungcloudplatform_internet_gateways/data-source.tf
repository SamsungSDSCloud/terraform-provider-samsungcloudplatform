data "samsungcloudplatform_internet_gateways" "pjt_igws" {
}

output "contents" {
  value = data.samsungcloudplatform_internet_gateways.pjt_igws.contents
}
