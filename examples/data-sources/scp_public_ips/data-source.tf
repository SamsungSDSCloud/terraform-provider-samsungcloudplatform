data "scp_region" "region" {
}

# Find public ip list
data "scp_public_ips" "my_scp_public_ips" {
  service_zone_id = data.scp_region.region.id
}

output "contents" {
  value = data.scp_public_ips.my_scp_public_ips.contents
}
