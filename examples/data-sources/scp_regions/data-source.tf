# Find all region list for current project
data "scp_regions" "my_regions1" {
}

output "contents" {
  value = data.scp_regions.my_regions1.regions
}
