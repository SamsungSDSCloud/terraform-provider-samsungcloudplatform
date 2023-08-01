data "scp_region" "region" {
}

resource "scp_vpc" "vpc4natgw" {
  name        = var.vpc_name
  description = "VPC for nat gateway"
  region      = data.scp_region.region.location
}

resource "scp_internet_gateway" "my_igw" {
  vpc_id      = scp_vpc.vpc4natgw.id
  description = "Internet GW generated from Terraform"
}

resource "scp_subnet" "subnet4natgw" {
  vpc_id      = scp_vpc.vpc4natgw.id
  name        = var.subnet_name
  type        = "PUBLIC"
  cidr_ipv4   = "192.169.4.0/24"
  description = "Subnet for nat gateway"
  depends_on  = [scp_internet_gateway.my_igw]
}

resource "scp_nat_gateway" "my_nat" {
  subnet_id   = scp_subnet.subnet4natgw.id
  description = "NAT GW from Terraform"
}
