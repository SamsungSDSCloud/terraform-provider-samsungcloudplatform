output "service_l4_id" {
  value = scp_lb_service.my_lb_service_l4.id
}

output "service_l7_id" {
  value = scp_lb_service.my_lb_service_l7.id
}
