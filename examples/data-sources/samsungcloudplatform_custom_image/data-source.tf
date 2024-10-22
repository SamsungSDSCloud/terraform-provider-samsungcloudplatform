data "samsungcloudplatform_region" "region" {
}

data "samsungcloudplatform_custom_images" "my_scp_custom_images" {
  service_group = "COMPUTE"
  service       = "Virtual Server"
  region = data.samsungcloudplatform_region.region.location
}

output "output_my_scp_custom_images" {
  value = data.samsungcloudplatform_custom_images.my_scp_custom_images
}
