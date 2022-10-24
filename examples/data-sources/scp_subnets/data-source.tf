data "scp_subnets" "my_scp_subnets1" {
}

output "output_my_scp_subnets1" {
  value = data.scp_subnets.my_scp_subnets1
}
