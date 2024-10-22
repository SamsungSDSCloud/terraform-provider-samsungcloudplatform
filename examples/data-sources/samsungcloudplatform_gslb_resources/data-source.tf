data "samsungcloudplatform_gslbs" "gslbs" {
}

data "samsungcloudplatform_gslb_resources" "my_scp_gslb_resources" {
  gslb_id = data.samsungcloudplatform_gslbs.gslbs.contents[0].gslb_id
}

output "contents" {
  value = data.samsungcloudplatform_gslb_resources.my_scp_gslb_resources.contents
}
