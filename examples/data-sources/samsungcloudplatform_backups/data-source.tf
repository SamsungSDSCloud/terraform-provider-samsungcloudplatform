data "samsungcloudplatform_backups" "my_scp_backups" {
  filter {
    name   = "backup_policy_type_category"
    values = ["VM"]
  }
}

output "output_my_scp_backups" {
  value = data.samsungcloudplatform_backups.my_scp_backups
}
