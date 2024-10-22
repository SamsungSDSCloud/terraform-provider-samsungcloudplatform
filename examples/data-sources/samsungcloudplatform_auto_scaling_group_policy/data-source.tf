# Find details of Auto-Scaling Group policy
data "samsungcloudplatform_auto_scaling_group_policy" "my_auto_scaling_group_policy" {
  asg_id = "AUTO_SCALING_GROUP-XXXXX"
  policy_id = "ASG_POLICY-XXXXX"
}

output "output_scp_auto_scaling_group_policy1" {
  value = data.samsungcloudplatform_auto_scaling_group_policy.my_auto_scaling_group_policy
}
