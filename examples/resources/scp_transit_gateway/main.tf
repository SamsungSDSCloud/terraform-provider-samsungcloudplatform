data "scp_region" "my_region" {
}

resource "scp_transit_gateway" "tgw01" {
  transit_gateway_name  = var.name
  transit_gateway_description = "create transit gateway from Terraform"
  region = data.scp_region.my_region.location
  bandwidth_gbps   = var.bandwidthGbps
  uplink_enabled   = var.uplinkEnabled
}

