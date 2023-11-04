# Find public ip
data "scp_public_ip" "my_public_ip" {
  public_ip_id = "PUBLIC_IP-XXXXXXXXXX"
}

output "output_my_public_ip" {
  value = data.scp_public_ip.my_public_ip
}
