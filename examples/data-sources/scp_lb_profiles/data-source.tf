data "scp_lb_profiles" "my_scp_lb_profiles" {
  load_balancer_id = "Load balancer id"
}
output "output_my_scp_lb_profiles" {
  value = data.scp_lb_profiles.my_scp_lb_profiles
}
