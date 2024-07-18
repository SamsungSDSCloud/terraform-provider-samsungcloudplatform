data "scp_region" "region" {
}

data "scp_custom_images" "my_scp_custom_images" {
  service_group = "COMPUTE"
  service       = "Virtual Server"
  region = data.scp_region.region.location
}

output "output_my_scp_custom_images" {
  value = data.scp_custom_images.my_scp_custom_images
}
