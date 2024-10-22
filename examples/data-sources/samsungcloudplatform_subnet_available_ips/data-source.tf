data "samsungcloudplatform_subnets" "my_scp_subnets" {
}

data "samsungcloudplatform_subnet_available_ips" "my_scp_subnet_available_ips1" {
  subnet_id  = data.samsungcloudplatform_subnets.my_scp_subnets.contents[0].subnet_id
}

output "contents" {
  value = data.samsungcloudplatform_subnet_available_ips.my_scp_subnet_available_ips1
}
