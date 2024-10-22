# Find all region list for current project
data "samsungcloudplatform_regions" "my_regions1" {
}

output "contents" {
  value = data.samsungcloudplatform_regions.my_regions1.regions
}
