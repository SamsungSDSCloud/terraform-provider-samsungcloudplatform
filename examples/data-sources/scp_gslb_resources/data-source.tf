data "scp_gslbs" "gslbs" {
}

data "scp_gslb_resources" "my_scp_gslb_resources" {
  gslb_id = data.scp_gslbs.gslbs.contents[0].gslb_id
}

output "contents" {
  value = data.scp_gslb_resources.my_scp_gslb_resources.contents
}
