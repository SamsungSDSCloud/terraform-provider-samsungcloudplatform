data "samsungcloudplatform_project" "my_scp_project" {
}

output "output_my_scp_project" {
  value = data.samsungcloudplatform_project.my_scp_project
}


