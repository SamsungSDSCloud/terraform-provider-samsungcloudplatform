data "scp_region" "region" {
}

resource "scp_vpc" "vpc4fw" {
  name        = var.name
  description = "VPC for firewall"
  region      = data.scp_region.region.location
}

resource "scp_internet_gateway" "vpc4fw_igw" {
  vpc_id      = scp_vpc.vpc4fw.id
  description = "IGW for VPC FW"
}

resource "scp_firewall" "vpc4fw_fw" {
  vpc_id    = scp_vpc.vpc4fw.id
  target_id = scp_internet_gateway.vpc4fw_igw.id
  enabled   = var.enabled
}
