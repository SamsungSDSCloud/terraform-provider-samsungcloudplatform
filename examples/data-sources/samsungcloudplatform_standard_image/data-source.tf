data "samsungcloudplatform_region" "region" {
}

# Find standard image
data "samsungcloudplatform_standard_image" "ubuntu_image" {
  service_group = "COMPUTE"
  service       = "Virtual Server"
  region        = data.samsungcloudplatform_region.region.location

  # Apply filter for 'image_name' regex value "Ubuntu 18.04 *"
  filter {
    name      = "image_name"
    values    = ["Ubuntu 18.04 *"]
    use_regex = true
  }
}

output "result_scp_my_standard_image" {
  value = data.samsungcloudplatform_standard_image.ubuntu_image
}
