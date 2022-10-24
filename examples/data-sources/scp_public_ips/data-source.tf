data "scp_region" "region" {
}

# Find public ip list
data "scp_public_ips" "my_scp_public_ips" {
  service_zone_id = data.scp_region.region.id
}

output "output_scp_public_ips" {
  value = data.scp_public_ips.my_scp_public_ips
}
