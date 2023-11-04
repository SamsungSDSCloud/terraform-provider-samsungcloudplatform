data "scp_resources_in_my_project" "my_resources" {

}

output "scp_resources_in_my_project" {
  value = data.scp_resources_in_my_project.my_resources
}
