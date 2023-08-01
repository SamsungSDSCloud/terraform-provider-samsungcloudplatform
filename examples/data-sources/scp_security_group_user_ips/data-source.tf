data "scp_security_group" "my_sg" {
  security_group_id = "FIREWALL_SECURITY_GROUP-XXXXXXXXXXXXXXXXXXXXXX"
}

output "output_my_scp_sg" {
  value = data.scp_security_group.my_sg.security_group_id
}

data "scp_security_group_user_ips" "my_sg_user_ips" {
  security_group_id = data.scp_security_group.my_sg.security_group_id
}

output "output_my_scp_sg_user_ips" {
  value = data.scp_security_group_user_ips.my_sg_user_ips
}
