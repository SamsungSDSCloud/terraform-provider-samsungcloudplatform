data "samsungcloudplatform_resource_group_resources" "my_resource_group_resources" {
  resource_group_id = "RESOURCE_GROUP-XXXXXXXXXXXXX"
}

output "result_my_resource_group_resources" {
  value = data.samsungcloudplatform_resource_group_resources.my_resource_group_resources
}
