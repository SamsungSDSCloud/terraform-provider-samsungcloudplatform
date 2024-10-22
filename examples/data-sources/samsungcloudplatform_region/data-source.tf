# Find region for current project
data "samsungcloudplatform_region" "my_region" {
}

output "output_my_scp_region" {
  value = data.samsungcloudplatform_region.my_region
}
