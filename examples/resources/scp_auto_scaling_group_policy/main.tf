# create two policies
resource "scp_auto_scaling_group_policy" "my_asg_policy1" {
  asg_id = data.terraform_remote_state.asg.outputs.id
  comparison_operator = var.comparison_operator
  cooldown_seconds = var.cooldown_seconds
  evaluation_minutes       = var.evaluation_minutes
  metric_method  = var.metric_method
  metric_type  = var.metric_type
  policy_name = var.name
  scale_method = var.scale_method
  scale_type = var.scale_type
  scale_value = var.scale_value
  threshold = var.threshold
}

resource "scp_auto_scaling_group_policy" "my_asg_policy2" {
  asg_id = data.terraform_remote_state.asg.outputs.id
  comparison_operator = var.comparison_operator_2
  cooldown_seconds = var.cooldown_seconds
  evaluation_minutes       = var.evaluation_minutes
  metric_method  = var.metric_method
  metric_type  = var.metric_type
  policy_name = var.name_2
  scale_method = var.scale_method
  scale_type = var.scale_type_2
  scale_value = var.scale_value
  threshold = var.threshold_2
}
