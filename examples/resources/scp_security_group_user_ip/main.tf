data "scp_region" "region" {
}

resource "scp_security_group_user_ip" "tc_sg_user_ip" {
  security_group_id = data.terraform_remote_state.security-group.outputs.id
  user_ip_type = "MULTICAST"
  user_ip_address = "224.0.0.1"
  user_ip_description = "Made in Terraform"
}
