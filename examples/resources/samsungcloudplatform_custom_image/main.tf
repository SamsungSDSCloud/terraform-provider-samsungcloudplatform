resource "samsungcloudplatform_custom_image" "custom_image_001" {
  image_name = var.name
  image_description = var.desc
  origin_virtual_server_id = data.samsungcloudplatform_virtual_servers.virtual_server_list.contents[0].virtual_server_id
}
