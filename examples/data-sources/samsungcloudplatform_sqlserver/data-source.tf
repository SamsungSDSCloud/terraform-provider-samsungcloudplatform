data "samsungcloudplatform_sqlserver" "my_scp_sqlserver" {
  sqlserver_cluster_id = "SERVICE-123456789"
}

output "output_my_scp_sqlserver" {
  value = data.samsungcloudplatform_sqlserver.my_scp_sqlserver
}
