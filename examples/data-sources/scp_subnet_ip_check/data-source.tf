data "scp_subnets" "my_scp_subnets" {
}

data "scp_subnet_ip_check" "scp_subnet_ip_check1" {
  subnet_id  = data.scp_subnets.my_scp_subnets.contents[0].subnet_id
  ip_address = "192.169.3.2"
}

output "result_scp_subnet_ip_check1" {
  value = data.scp_subnet_ip_check.scp_subnet_ip_check1
}
