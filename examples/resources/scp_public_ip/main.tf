data "scp_region" "region" {
}

resource "scp_public_ip" "ip01" {
  description = "Public IP generated from Terraform"
  region      = data.scp_region.region.location
}
