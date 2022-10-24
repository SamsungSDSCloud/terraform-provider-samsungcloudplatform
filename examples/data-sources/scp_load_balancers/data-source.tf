data "scp_load_balancers" "my_scp_load_balancers" {
}

output "result_scp_load_balancers" {
  value = data.scp_load_balancers.my_scp_load_balancers
}
