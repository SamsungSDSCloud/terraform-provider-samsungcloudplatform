resource "samsungcloudplatform_subnet_security_group" "my_subnet_security_group" {
  subnet_id            = "SUBNET-xxxxxxxxxxxxxxx"
  vip_id               = "SUBNET_VIRTUAL_IP-xxxxxxxxxxxxxxxp"
  security_group_id = "FIREWALL_SECURITY_GROUP-xxxxxxxxxxxxxxx"
}
