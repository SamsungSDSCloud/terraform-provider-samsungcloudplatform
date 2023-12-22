resource "scp_auto_scaling_group_load_balancer" "my_auto_scaling_group_load_balancer" {
  asg_id = data.terraform_remote_state.auto_scaling_group.outputs.id
  lb_rule_ids = var.lb_rule_ids
}
