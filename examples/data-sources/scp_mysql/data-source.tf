data "scp_mysql" "my_scp_mysql" {
  mysql_cluster_id = "SERVICE-123456789"
}

output "output_my_scp_mysql" {
  value = data.scp_mysql.my_scp_mysql
}
