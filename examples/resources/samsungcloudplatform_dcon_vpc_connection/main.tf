data "samsungcloudplatform_region" "region" {
}

resource "samsungcloudplatform_vpc" "connect_vpc" {
  name        = var.vpc_name
  description = "VPC for connection"
  region      = data.samsungcloudplatform_region.region.location
}

resource "samsungcloudplatform_direct_connect" "connect_dc" {
  name        = var.dc_name
  description = "dc for connection"
  region      = data.samsungcloudplatform_region.region.location
  bandwidth   = var.bandwidth
}

resource "samsungcloudplatform_dcon_vpc_connection" "dconvpc01" {
  vpc_id            = samsungcloudplatform_vpc.connect_vpc.id
  firewall_enabled  = var.firewall_enabled
  direct_connect_id = samsungcloudplatform_direct_connect.connect_dc.id
  description       = "Dcon-vpc connection generated from Terraform"
}

