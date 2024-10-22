data "samsungcloudplatform_security_group" "my_sg" {
  security_group_id = "FIREWALL_SECURITY_GROUP-XXXXXXXXXXXXXXXXXXXXXX"
}

output "output_my_scp_sg" {
  value = data.samsungcloudplatform_security_group.my_sg
}
