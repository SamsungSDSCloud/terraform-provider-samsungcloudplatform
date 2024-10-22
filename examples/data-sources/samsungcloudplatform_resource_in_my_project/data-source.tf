data "samsungcloudplatform_resource_in_my_project" "my_resource" {
  resource_id = "RESOURCE-XXXXX"
  project_id = "PROJECT-XXXXXXXXXXXX"
}

output "result_my_resource" {
  value = data.samsungcloudplatform_resource_in_my_project.my_resource
}
