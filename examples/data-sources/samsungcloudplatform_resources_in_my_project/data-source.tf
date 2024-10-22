data "samsungcloudplatform_resources_in_my_project" "my_resources" {

}

output "samsungcloudplatform_resources_in_my_project" {
  value = data.samsungcloudplatform_resources_in_my_project.my_resources
}
