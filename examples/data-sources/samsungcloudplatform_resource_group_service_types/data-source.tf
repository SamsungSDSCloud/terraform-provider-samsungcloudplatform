data "samsungcloudplatform_resource_group_service_types" "my_service_types" {
}

output "result_my_service_types" {
  value = data.samsungcloudplatform_resource_group_service_types.my_service_types
}
