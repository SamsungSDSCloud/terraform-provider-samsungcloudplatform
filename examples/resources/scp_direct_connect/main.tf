data "scp_region" "my_region" {
}

resource "scp_direct_connect" "dc01" {
  name        = var.name
  description = "DirectConnect generated from Terraform"
  region      = data.scp_region.my_region.location
  bandwidth   = var.bandwidth
}

