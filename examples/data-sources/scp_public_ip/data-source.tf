data "scp_region" "my_region" {
}

# Find public ip
data "scp_public_ip" "my_public_ip" {
  region = data.scp_region.my_region.location
}

output "output_my_public_ip" {
  value = data.scp_public_ip.my_public_ip
}
