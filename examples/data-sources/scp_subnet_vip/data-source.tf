data "scp_subnet_vip_detail" "my_scp_subnet_vip1" {
  subnet_id = "SUBNET-XXXXXXXXXXXX"
  vip_id = "SUBNET_VIRTUAL_IP-XXXXXXXXXXXX"
}

output "output_my_scp_subnet_vip1" {
  value = data.scp_subnet_vip_detail.my_scp_subnet_vip1
}
