data "scp_resource_groups" "my_resource_groups" {
}

output "result_my_resource_groups" {
  value = data.scp_resource_groups.my_resource_groups
}
