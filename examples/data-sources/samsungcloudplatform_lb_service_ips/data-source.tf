data "samsungcloudplatform_lb_service_ips" "my_scp_lb_service_ips" {
  load_balancer_id = "lb id"
}

output "output_my_scp_lb_service_ips" {
  value = data.samsungcloudplatform_lb_service_ips.my_scp_lb_service_ips
}
