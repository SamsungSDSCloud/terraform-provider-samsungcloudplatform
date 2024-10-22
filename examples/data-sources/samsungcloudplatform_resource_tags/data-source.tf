data "samsungcloudplatform_resource_tags" "my_resource_tags" {
  resource_id = "RESOURCE_XXXXXXXX"
}

output "contents" {
  value = data.samsungcloudplatform_resource_tags.my_resource_tags.contents
}
