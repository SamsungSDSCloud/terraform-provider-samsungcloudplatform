data "samsungcloudplatform_direct_connects" "pjt_dcs" {

}

output "contents" {
  value = data.samsungcloudplatform_direct_connects.pjt_dcs.contents
}
