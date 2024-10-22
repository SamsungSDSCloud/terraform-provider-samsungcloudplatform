data "samsungcloudplatform_resource_groups" "my_resource_groups" {
}

output "result_my_resource_groups" {
  value = data.samsungcloudplatform_resource_groups.my_resource_groups
}
