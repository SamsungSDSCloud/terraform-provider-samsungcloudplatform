data "scp_backups" "my_scp_backups" {
  filter {
    name   = "backup_policy_type_category"
    values = ["VM"]
  }
}

output "output_my_scp_backups" {
  value = data.scp_backups.my_scp_backups
}
