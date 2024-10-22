data "samsungcloudplatform_region" "region" {
}

resource "samsungcloudplatform_vpc" "vpc4fw" {
  name        = var.name
  description = "VPC for firewall"
  region      = data.samsungcloudplatform_region.region.location
}

resource "samsungcloudplatform_internet_gateway" "vpc4fw_igw" {
  vpc_id      = samsungcloudplatform_vpc.vpc4fw.id
  description = "IGW for VPC FW"
}

resource "samsungcloudplatform_firewall" "vpc4fw_fw" {
  vpc_id    = samsungcloudplatform_vpc.vpc4fw.id
  target_id = samsungcloudplatform_internet_gateway.vpc4fw_igw.id
  enabled   = var.enabled
  logging_enabled = var.logging_enabled
}
