data "samsungcloudplatform_resource_group_in_my_projects" "my_resource_group_in_my_projects" {
  resource_group_id = "RESOURCE_GROUP-XXXXXXXXXXXXX"
}

output "result_my_resource_group" {
  value = data.samsungcloudplatform_resource_group_in_my_projects.my_resource_group_in_my_projects
}
