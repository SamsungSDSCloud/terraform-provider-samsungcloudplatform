data "samsungcloudplatform_security_group_rule" "my_sg_rule" {
  security_group_id = "FIREWALL_SECURITY_GROUP-XXXXXXXXXXXXXXXXXXXXXX"
  rule_id = "FIREWALL_RULE-XXXXXXXXXXXXXXXXXXXXXX"
}

output "output_my_scp_sg_rule" {
  value = data.samsungcloudplatform_security_group_rule.my_sg_rule
}
