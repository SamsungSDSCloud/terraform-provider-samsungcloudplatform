data "samsungcloudplatform_security_group_user_ips" "my_sg_user_ips" {
  security_group_id = "FIREWALL_SECURITY_GROUP-XXXXXXXXXXXXXXXXXXXXXX"
}

output "output_my_scp_sg_user_ips" {
  value = data.samsungcloudplatform_security_group_user_ips.my_sg_user_ips
}
