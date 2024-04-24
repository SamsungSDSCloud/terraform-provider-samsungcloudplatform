# Find all Auto-Scaling Group Virtual Servers
data "scp_auto_scaling_group_virtual_servers" "my_auto_scaling_group_virtual_servers" {
  asg_id = "AUTO_SCALING_GROUP-XXXXX"
}

output "output_scp_auto_scaling_group_virtual_servers" {
  value = data.scp_auto_scaling_group_virtual_servers.my_auto_scaling_group_virtual_servers
}
