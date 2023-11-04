data "scp_resource_group_service_types" "my_service_types" {
}

output "result_my_service_types" {
  value = data.scp_resource_group_service_types.my_service_types
}
