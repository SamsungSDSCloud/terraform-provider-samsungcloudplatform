data "scp_security_group_rules" "my_sg_rules" {
  security_group_id = "FIREWALL_SECURITY_GROUP-XXXXXXXXXXXXXXXXXXXXXX"
}

data "scp_security_group_rules" "my_sg_rules2" {
  security_group_id = data.scp_security_group_rules.my_sg_rules.security_group_id
  filter {
    name   = "rule_direction"
    values = ["IN"]
  }
}

output "output_my_scp_sg_rules" {
  value = data.scp_security_group_rules.my_sg_rules
}

output "output_my_scp_sg_rules2" {
  value = data.scp_security_group_rules.my_sg_rules2
}
