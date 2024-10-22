data "samsungcloudplatform_nat_gateways" "pjt_natgws" {

}

output "contents" {
  value = data.samsungcloudplatform_nat_gateways.pjt_natgws.contents
}
