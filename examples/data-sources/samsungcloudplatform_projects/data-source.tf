data "samsungcloudplatform_projects" "my_scp_projects" {
}

output "output_my_scp_project" {
  value = data.samsungcloudplatform_projects.my_scp_projects
}


