resource "scp_lb_service" "my_lb_service_l4" {
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
  public_ip_id     = scp_public_ip.my_public_ip_id.id
}

resource "scp_lb_service" "my_lb_service_l7" {
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

data "scp_region" "my_region" {
}

resource "scp_public_ip" "my_public_ip_id" {
  region = data.scp_region.my_region.location
}
