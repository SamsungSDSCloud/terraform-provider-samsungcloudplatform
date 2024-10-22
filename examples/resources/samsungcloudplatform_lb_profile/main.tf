resource "samsungcloudplatform_lb_profile" "my_lb_profile_persistence" {
  lb_id            = data.terraform_remote_state.load_balancer.outputs.id
  name             = var.name_persistence
  category         = "PERSISTENCE"
  persistence_type = "SOURCE_IP"
}

resource "samsungcloudplatform_lb_profile" "my_lb_profile_app_l4" {
  lb_id           = data.terraform_remote_state.load_balancer.outputs.id
  name            = var.name
  category        = "APPLICATION"
  layer_type      = "L4"
  session_timeout = 30
}

resource "samsungcloudplatform_lb_profile" "my_lb_profile_app_l7" {
  lb_id                = data.terraform_remote_state.load_balancer.outputs.id
  name                 = var.name_l7
  category             = "APPLICATION"
  layer_type           = "L7"
  redirect_type        = "HTTP_TO_HTTPS_REDIRECT"
  request_header_size  = 1
  response_header_size = 1
  response_timeout     = 1
  session_timeout      = 30
  x_forwarded_for      = "REPLACE"
}
