data "samsungcloudplatform_resource_group_resource_types" "my_resource_types" {

}

output "result_my_resource_types" {
  value = data.samsungcloudplatform_resource_group_resource_types.my_resource_types
}
