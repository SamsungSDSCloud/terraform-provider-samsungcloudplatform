data "scp_resource_group_resource_types" "my_resource_types" {

}

output "result_my_resource_types" {
  value = data.scp_resource_group_resource_types.my_resource_types
}
