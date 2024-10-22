data "samsungcloudplatform_resource_groups_in_my_projects" "my_resource_groups_in_my_projects" {
}

output "result_my_resource_groups" {
  value = data.samsungcloudplatform_resource_groups_in_my_projects.my_resource_groups_in_my_projects
}
