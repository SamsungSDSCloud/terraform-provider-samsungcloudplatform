resource "scp_custom_image" "custom_image_001" {
  image_name = var.name
  image_description = var.desc
  origin_virtual_server_id = data.virtual_server_list.contents[0].virtual_server_id
}
