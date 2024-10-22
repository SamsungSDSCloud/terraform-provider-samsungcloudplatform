data "samsungcloudplatform_resources" "my_resources" {
}

output "out_resource" {
  value = data.samsungcloudplatform_resources.my_resources
}
