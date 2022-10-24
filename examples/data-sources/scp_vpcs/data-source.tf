# Find all vpcs for current project
data "scp_vpcs" "my_scp_vpcs1" {
}

output "output_my_scp_vpcs1" {
  value = data.scp_vpcs.my_scp_vpcs1
}
