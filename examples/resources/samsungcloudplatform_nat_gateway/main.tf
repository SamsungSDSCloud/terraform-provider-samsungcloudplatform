data "samsungcloudplatform_region" "region" {
}

resource "samsungcloudplatform_vpc" "vpc4natgw" {
  name        = var.vpc_name
  description = "VPC for nat gateway"
  region      = data.samsungcloudplatform_region.region.location
}

resource "samsungcloudplatform_internet_gateway" "my_igw" {
  vpc_id      = samsungcloudplatform_vpc.vpc4natgw.id
  description = "Internet GW generated from Terraform"
}

resource "samsungcloudplatform_subnet" "subnet4natgw" {
  vpc_id      = samsungcloudplatform_vpc.vpc4natgw.id
  name        = var.subnet_name
  type        = "PUBLIC"
  cidr_ipv4   = "192.169.4.0/24"
  description = "Subnet for nat gateway"
  depends_on  = [samsungcloudplatform_internet_gateway.my_igw]
}

resource "samsungcloudplatform_nat_gateway" "my_nat" {
  subnet_id   = samsungcloudplatform_subnet.subnet4natgw.id
  description = "NAT GW from Terraform"
}
