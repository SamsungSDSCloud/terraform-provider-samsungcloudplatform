data "scp_project" "my_scp_project" {
}

output "output_my_scp_project" {
  value = data.scp_project.my_scp_project
}


