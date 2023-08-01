data "scp_subnets" "my_scp_subnets" {
}

data "scp_subnet_available_ips" "my_scp_subnet_available_ips1" {
  subnet_id  = data.scp_subnets.my_scp_subnets.contents[0].subnet_id
}

output "contents" {
  value = data.scp_subnet_available_ips.my_scp_subnet_available_ips1
}
