data "scp_region" "region" {
}

resource "scp_security_group" "tc_sg" {
  vpc_id      = data.terraform_remote_state.vpc.outputs.id
  name        = var.name
  description = "SecurityGroup generated from terraform"
  is_loggable = false
}

resource "scp_security_group_user_ip" "tc_sg_user_ip" {
  security_group_id = scp_security_group.tc_sg.id
  user_ip_type = "MULTICAST"
  user_ip_address = "224.0.0.0"
  user_ip_description = "Made in Terraform"
}

resource "scp_security_group_rule" "tc_rule_all" {
  security_group_id = scp_security_group.tc_sg.id
  direction         = "in"
  description       = "SecurityGroup Rule generated from Terraform"
  addresses_ipv4 = [
    "0.0.0.0/0"
  ]
  service {
    type = "all"
  }
}
