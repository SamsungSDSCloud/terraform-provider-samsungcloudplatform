data "scp_subnets" "my_scp_subnets" {
}

output "contents" {
  value = data.scp_subnets.my_scp_subnets.contents
}
