# Find details of Launch Configuration
data "samsungcloudplatform_auto_scaling_group" "my_auto_scaling_group1" {
  asg_id = "AUTO_SCALING_GROUP-XXXXX"
}

output "output_scp_auto_scaling_group1" {
  value = data.samsungcloudplatform_auto_scaling_group.my_auto_scaling_group1
}
