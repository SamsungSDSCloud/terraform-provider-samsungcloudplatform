data "scp_epas" "my_scp_epas" {
  epas_cluster_id = "SERVICE-123456789"
}

output "output_my_scp_epas" {
  value = data.scp_epas.my_scp_epas
}
