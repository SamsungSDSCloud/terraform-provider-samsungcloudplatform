data "scp_backups" "my_scp_backups" {
}

output "output_my_scp_backups" {
  value = data.scp_backups.my_scp_backups
}
