data "scp_resource_group" "my_resource_group" {
  resource_group_id = "RESOURCE_GROUP-XXXXXXXXXXXXX"
}

output "result_my_resource_group" {
  value = data.scp_resource_group.my_resource_group
}
