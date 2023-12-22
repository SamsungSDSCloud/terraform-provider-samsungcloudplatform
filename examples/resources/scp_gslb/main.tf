resource "scp_gslb" "my_gslb" {
  gslb_name        = var.name
  gslb_env_usage = "PUBLIC"
  gslb_algorithm = "RATIO"
  protocol = "HTTP"
  gslb_health_check_interval = 14
  gslb_health_check_timeout = 15
  probe_timeout = 10
  service_port = 80
  gslb_health_check_user_id = "tester"
  gslb_health_check_user_password = "test123$%"
  gslb_send_string = "GET /index.html"
  gslb_response_string= "test response"
  gslb_resources  {
    gslb_destination            = "192.168.0.1"
    gslb_region    = "KR-WEST-1"
    gslb_resource_weight = 10
    gslb_resource_description       = "test resource 1"
  }
  gslb_resources {
      gslb_destination            = "192.168.0.2"
      gslb_region    = "KR-EAST-1"
      gslb_resource_weight = 20
      gslb_resource_description       = "test resource 2"
  }
  gslb_resources {
    gslb_destination            = "192.168.0.3"
    gslb_region    = "KR-EAST-1"
    gslb_resource_weight = 15
    gslb_resource_description       = "test resource 3"
  }
}
#
#
#resource "scp_gslb" "gslbtftest02" {
#  gslb_name        = "gslbtftest02"
#  gslb_env_usage = "PUBLIC"
#  gslb_algorithm = "RTT"
#  protocol = "HTTPS"
#  gslb_health_check_interval = 14
#  gslb_health_check_timeout = 15
#  probe_timeout = 10
#  service_port = 80
#  gslb_health_check_user_id = "tester"
#  gslb_health_check_user_password = "test123$%"
#  gslb_send_string = "GET /index.html"
#  gslb_response_string= "test response"
#  gslb_resources  {
#    gslb_destination            = "192.168.0.1"
#    gslb_region    = "KR-WEST-1"
#    gslb_resource_weight = 10
#    gslb_resource_description       = "test resource 1"
#  }
#  gslb_resources {
#    gslb_destination            = "192.168.0.3"
#    gslb_region    = "KR-EAST-1"
#    gslb_resource_weight = 20
#    gslb_resource_description       = "test resource 2"
#  }
#}
#
#resource "scp_gslb" "gslbtftest03" {
#  gslb_name        = "gslbtftest03"
#  gslb_env_usage = "PUBLIC"
#  gslb_algorithm = "RATIO"
#  protocol = "TCP"
#  gslb_health_check_interval = 10
#  gslb_health_check_timeout = 15
#  probe_timeout = 10
#  service_port = 22
#  gslb_resources  {
#    gslb_destination            = "192.168.0.1"
#    gslb_region    = "KR-WEST-1"
#    gslb_resource_weight = 20
#    gslb_resource_description       = "test resource 1"
#  }
#}
