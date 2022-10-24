resource "scp_lb_server_group" "my_lb_server_group_http" {
  lb_id     = data.terraform_remote_state.load_balancer.outputs.id
  name      = "${var.name}-http"
  algorithm = "ROUND_ROBIN"

  server_group_member {
    join_state  = "ENABLED"
    object_type = "INSTANCE"
    object_id   = data.terraform_remote_state.virtual_server.outputs.id
    object_port = 8116
    weight      = 2
  }

  monitor_protocol     = "HTTP"
  monitor_port         = 8020
  monitor_interval_sec = 30
  monitor_timeout_sec  = 30
  monitor_count        = 30

  monitor_http_method        = "POST"
  monitor_http_version       = "1.1"
  monitor_http_url           = "/health-check-post"
  monitor_http_request_body  = "TestPost1"
  monitor_http_response_body = "TestPost2"
}

resource "scp_lb_server_group" "my_lb_server_group_tcp" {
  lb_id     = "LB-xxxxxxx"
  name      = "${var.name}-tcp"
  algorithm = "LEAST_CONNECTION"

  server_group_member {
    join_state  = "ENABLED"
    object_type = "INSTANCE"
    object_id   = data.terraform_remote_state.virtual_server.outputs.id
    object_port = 8117
    weight      = 10
  }

  monitor_protocol     = "TCP"
  monitor_port         = 8020
  monitor_interval_sec = 30
  monitor_timeout_sec  = 30
  monitor_count        = 30
}
