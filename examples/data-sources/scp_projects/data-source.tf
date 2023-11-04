data "scp_projects" "my_scp_projects" {
}

output "output_my_scp_project" {
  value = data.scp_projects.my_scp_projects
}


