data "samsungcloudplatform_region" "region" {
}

resource "samsungcloudplatform_vpc" "vpc4igw" {
  name        = var.name
  description = "VPC for internet gateway"
  region      = data.samsungcloudplatform_region.region.location
}

resource "samsungcloudplatform_internet_gateway" "my_igw" {
  vpc_id      = samsungcloudplatform_vpc.vpc4igw.id
  igw_type    = var.type
  description = "Internet GW generated from Terraform"
}
