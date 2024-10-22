data "samsungcloudplatform_resource_group_resources_in_my_projects" "my_resource_group_resources_in_my_projects" {
  resource_group_id = "RESOURCE_GROUP-XXXXXXXXXXXXX"
}

output "result_my_resource_group_resources" {
  value = data.samsungcloudplatform_resource_group_resources_in_my_projects.my_resource_group_resources_in_my_projects
}
