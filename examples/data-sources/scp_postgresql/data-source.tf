data "scp_postgresql" "my_scp_postgresql" {
  postgresql_cluster_id = "SERVICE-123456789"
}

output "output_my_scp_postgresql" {
  value = data.scp_postgresql.my_scp_postgresql
}
