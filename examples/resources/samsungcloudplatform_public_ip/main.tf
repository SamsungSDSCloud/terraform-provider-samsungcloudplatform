data "samsungcloudplatform_region" "region" {
}

resource "samsungcloudplatform_public_ip" "ip01" {
  description = "Public IP generated from Terraform"
  region      = data.samsungcloudplatform_region.region.location
  uplink_type = "INTERNET"
}
