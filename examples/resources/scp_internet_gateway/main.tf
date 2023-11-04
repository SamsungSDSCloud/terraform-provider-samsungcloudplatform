data "scp_region" "region" {
}

resource "scp_vpc" "vpc4igw" {
  name        = var.name
  description = "VPC for internet gateway"
  region      = data.scp_region.region.location
}

resource "scp_internet_gateway" "my_igw" {
  vpc_id      = scp_vpc.vpc4igw.id
  igw_type    = var.type
  description = "Internet GW generated from Terraform"
}
