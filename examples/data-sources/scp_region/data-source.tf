# Find region for current project
data "scp_region" "my_region" {
}

output "output_my_scp_region" {
  value = data.scp_region.my_region
}
