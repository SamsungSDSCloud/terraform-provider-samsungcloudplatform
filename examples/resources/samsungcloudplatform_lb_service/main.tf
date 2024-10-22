resource "samsungcloudplatform_lb_service" "my_lb_service_l4" {
  lb_id            = data.terraform_remote_state.load_balancer.outputs.id
  name             = var.namel4
  layer_type       = "L4"
  protocol         = "TCP"
  service_ports    = "8090"
  forwarding_ports = "8091"
  service_ipv4     = "192.168.102.20"
  persistence      = "SOURCE_IP"
  app_profile_id   = data.terraform_remote_state.load_balancer_profile.outputs.id
  persistence_profile_id = data.terraform_remote_state.load_balancer_profile.outputs.persistence_id
  nat_active       = true
  public_ip_id     = samsungcloudplatform_public_ip.my_public_ip_id.id
}

resource "samsungcloudplatform_lb_service" "my_lb_service_l7" {
  lb_id            = data.terraform_remote_state.load_balancer.outputs.id
  app_profile_id   = data.terraform_remote_state.load_balancer_profile.outputs.id
  name             = var.namel7
  layer_type       = "L7"
  protocol         = "HTTP"
  service_ports    = "8088"
  forwarding_ports = "8089"
  service_ipv4     = "192.168.102.11"
  persistence      = "DISABLED"

  lb_rules {
    lb_rule_seq = 1
    pattern_url = "/promise"
  }
  lb_rules {
    lb_rule_seq = 2
    pattern_url = "/devotion"
  }
  nat_active    = false
}

resource "samsungcloudplatform_lb_service" "my_lb_service_l7_https" {
  lb_id            = data.terraform_remote_state.load_balancer.outputs.id
  app_profile_id   = data.terraform_remote_state.load_balancer_profile.outputs.id
  name             = var.namel7https
  layer_type       = "L7"
  protocol         = "HTTPS"
  client_certificate_id = "CERT-XXXXXXXXXXXXXXXXXXXXXX"
  client_ssl_security_level = "HIGH"
  server_certificate_id = "CERT-XXXXXXXXXXXXXXXXXXXXXX"
  server_ssl_security_level = "HIGH"
  service_ports    = "8088"
  forwarding_ports = "8089"
  service_ipv4     = "192.168.102.13"
  persistence      = "DISABLED"

  lb_rules {
    lb_rule_seq = 1
    pattern_url = "/promise"
  }
  lb_rules {
    lb_rule_seq = 2
    pattern_url = "/devotion"
  }
  nat_active    = false
}


data "samsungcloudplatform_region" "my_region" {
}

resource "samsungcloudplatform_public_ip" "my_public_ip_id" {
  region = data.samsungcloudplatform_region.my_region.location
  uplink_type = "INTERNET"
}
