data "scp_region" "region" {
}

resource "scp_security_group" "tc_sg" {
  vpc_id      = data.terraform_remote_state.vpc.outputs.id
  name        = var.name
  description    = "SG Policy for Bulk Rule generated from Terraform"
  is_loggable = false
}

resource "scp_security_group_bulk_rule" "tc_rule_all" {
  security_group_id = scp_security_group.tc_sg.id
  rule {
    direction      = "in"
    description    = "SecurityGroup Bulk Rule generated from Terraform"
    addresses_ipv4 = [
      "10.10.10.10", "20.20.20.20"
    ]
    service {
      type  = "tcp"
      value = "8080"
    }
    service {
      type  = "udp"
      value = "443"
    }
  }
  rule {
    direction      = "out"
    description    = "SecurityGroup Bulk Rule generated from Terraform"
    addresses_ipv4 = [
      "192.168.0.0/24"
    ]
    service {
      type = "all"
    }
  }
  rule {
    direction      = "in"
    description    = "SecurityGroup Bulk Rule generated from Terraform"
    addresses_ipv4 = [
      "30.30.30.30"
    ]
    service {
      type = "all"
    }
  }
}
