data "samsungcloudplatform_region" "my_region" {
}

resource "samsungcloudplatform_direct_connect" "dc01" {
  name        = var.name
  description = "DirectConnect generated from Terraform"
  region      = data.samsungcloudplatform_region.my_region.location
  bandwidth   = var.bandwidth
}

