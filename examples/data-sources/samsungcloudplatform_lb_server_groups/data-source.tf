# Find all load balancer server groups for current project
data "samsungcloudplatform_lb_server_groups" "my_scp_lb_server_groups" {
  load_balancer_id = "Load balancer id"
}

output "output_my_scp_lb_server_groups" {
  value = data.samsungcloudplatform_lb_server_groups.my_scp_lb_server_groups
}
