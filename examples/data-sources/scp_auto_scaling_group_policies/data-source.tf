# Find all Auto-Scaling Group policies
data "scp_auto_scaling_group_policies" "my_auto_scaling_group_policies1" {
  asg_id = "AUTO_SCALING_GROUP-XXXXX"
}

# Find all Auto-Scaling Group policies
data "scp_auto_scaling_group_policies" "my_auto_scaling_group_policies2" {
  asg_id = "AUTO_SCALING_GROUP-XXXXX"

  # Sort in ascending order of creation date
  sort = "createdDt:asc"

  # Apply filter for 'policy_name' regex value "test"
  filter {
    name = "policy_name"
    values = ["test"]
    use_regex = true
  }
}

output "output_scp_auto_scaling_group_policies1" {
  value = data.scp_auto_scaling_group_policies.my_auto_scaling_group_policies1
}

output "output_scp_auto_scaling_group_policies2" {
  value = data.scp_auto_scaling_group_policies.my_auto_scaling_group_policies2
}
