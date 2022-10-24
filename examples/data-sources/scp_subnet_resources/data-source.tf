data "scp_subnet_resources" "my_scp_subnet_resources1" {
  subnet_id = "SUBNET-xxxxx"
}

output "output_my_scp_subnet_resources1" {
  value = data.scp_subnet_resources.my_scp_subnet_resources1
}
