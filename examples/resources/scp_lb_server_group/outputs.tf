output "group_http_id" {
  value = scp_lb_server_group.my_lb_server_group_http.id
}

output "group_tcp_id" {
  value = scp_lb_server_group.my_lb_server_group_tcp.id
}
