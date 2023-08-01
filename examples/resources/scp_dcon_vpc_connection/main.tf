data "scp_region" "region" {
}

resource "scp_vpc" "connect_vpc" {
  name        = var.vpc_name
  description = "VPC for connection"
  region      = data.scp_region.region.location
}

resource "scp_direct_connect" "connect_dc" {
  name        = var.dc_name
  description = "dc for connection"
  region      = data.scp_region.region.location
  bandwidth   = var.bandwidth
}

resource "scp_dcon_vpc_connection" "dconvpc01" {
  vpc_id            = scp_vpc.connect_vpc.id
  firewall_enabled  = var.firewall_enabled
  direct_connect_id = scp_direct_connect.connect_dc.id
  description       = "Dcon-vpc connection generated from Terraform"
}

