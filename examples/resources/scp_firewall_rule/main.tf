data "scp_region" "region" {
}

resource "scp_firewall_rule" "vpc4fw_fwrule" {
  firewall_id = "FIREWALL-xxxxxxx"

  direction = "IN_OUT"
  action    = "ALLOW"

  enabled = false

  source_addresses_ipv4      = ["128.0.0.0/1"]
  destination_addresses_ipv4 = ["128.0.0.0/1"]

  service {
    type  = "TCP"
    value = "8080"
  }
  service {
    type  = "UDP"
    value = "22"
  }
  service {
    type = "TCP_ALL"
    value = ""
  }

  description = "Rule from terraform"
}
