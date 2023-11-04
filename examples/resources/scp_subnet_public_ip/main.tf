resource "scp_subnet_public_ip" "my_subnet_public_ip" {
  subnet_id            = "SUBNET-xxxxxxxxxxxxxxx"
  vip_id               = "SUBNET_VIRTUAL_IP-xxxxxxxxxxxxxxx"
  public_ip_address_id = "PUBLIC_IP-xxxxxxxxxxxxxxx"
}
