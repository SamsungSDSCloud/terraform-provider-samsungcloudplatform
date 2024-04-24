data "scp_custom_image" "my_scp_custom_image" {
  image_id = "IMAGE_XXXXX"
}

output "output_my_scp_custom_image" {
  value = data.scp_custom_image.my_scp_custom_image
}
