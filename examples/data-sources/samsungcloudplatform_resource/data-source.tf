data "samsungcloudplatform_resource" "my_resource" {
  resource_id = "VPC-XXXXXXX"
}

output "result_my_resource" {
  value = data.samsungcloudplatform_resource.my_resource
}
