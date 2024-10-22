data "samsungcloudplatform_lb_services" "my_scp_lb_services" {
  load_balancer_id = "lb id"
}
output "output_scp_public_ips" {
  value = data.samsungcloudplatform_lb_services.my_scp_lb_services
}
