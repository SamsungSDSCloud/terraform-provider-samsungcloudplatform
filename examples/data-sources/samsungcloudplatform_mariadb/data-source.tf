data "samsungcloudplatform_mariadb" "my_scp_mariadb" {
  mariadb_cluster_id = "SERVICE-123456789"
}

output "output_my_scp_mariadb" {
  value = data.samsungcloudplatform_mariadb.my_scp_mariadb
}
