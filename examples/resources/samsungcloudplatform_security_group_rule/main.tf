data "samsungcloudplatform_region" "region" {
}

resource "samsungcloudplatform_security_group_rule" "tc_my_rule" {
  security_group_id = data.terraform_remote_state.security-group.outputs.id
  direction         = "out"
  description       = "SecurityGroup Rule generated from Terraform"
  addresses_ipv4 = [
    "1.1.1.1/1"
  ]
  service {
    type  = "tcp"
    value = "80"
  }
}
