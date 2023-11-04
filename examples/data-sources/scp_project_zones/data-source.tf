data "scp_project_zones" "my_scp_project_zones" {
  project_id = "PROJECT-XXXXXXX"
}

output "output_my_scp_project" {
  value = data.scp_project_zones.my_scp_project_zones
}


