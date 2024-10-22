data "samsungcloudplatform_project_zones" "my_scp_project_zones" {
  project_id = "PROJECT-XXXXXXX"
}

output "output_my_scp_project" {
  value = data.samsungcloudplatform_project_zones.my_scp_project_zones
}


