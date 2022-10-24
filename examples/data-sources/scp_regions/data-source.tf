# Find all region list for current project
data "scp_regions" "my_regions1" {
}

output "output_my_scp_regions" {
  value = data.scp_regions.my_regions1
}
