data "samsungcloudplatform_region" "region" {
}

# Find all standard images
data "samsungcloudplatform_standard_images" "my_standard_images1" {
  service_group = "COMPUTE"
  service       = "Virtual Server"
  region        = data.samsungcloudplatform_region.region.location
}

# Find all standard images
data "samsungcloudplatform_standard_images" "my_standard_images2" {
  service_group = "COMPUTE"
  service       = "Virtual Server"
  region        = data.samsungcloudplatform_region.region.location

  # Apply filter for 'image_name' regex value "CentOS 7.8"
  filter {
    name   = "image_name"
    values = ["CentOS 7.8"]
  }
}

output "result_scp_my_standard_images1" {
  value = data.samsungcloudplatform_standard_images.my_standard_images1
}

output "result_scp_my_standard_images2" {
  value = data.samsungcloudplatform_standard_images.my_standard_images2
}
