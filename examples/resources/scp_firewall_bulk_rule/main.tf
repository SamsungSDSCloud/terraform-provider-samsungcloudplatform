data "scp_region" "region" {
}

resource "scp_firewall_bulk_rule" "vpc4fw_bulk_fwrule" {
  firewall_id = "FIREWALL-XXXX"

  bulk_rule_location_type = "LAST"
  bulk_rule_location_id = "FIREWALL_RULE-XXXX"

  rule {
    direction = "IN_OUT"
    action    = "ALLOW"

    enabled = false

    source_addresses_ipv4      = ["10.0.0.0/1"]
    destination_addresses_ipv4 = ["10.0.0.0/1"]

    service {
      type  = "TCP"
      value = "8080"
    }
    service {
      type  = "UDP"
      value = "22"
    }
    service {
      type  = "TCP_ALL"
      value = ""
    }

    description = "Bulk Rule 1 from terraform"
  }

  rule {
    direction = "IN"
    action    = "ALLOW"

    enabled = false

    source_addresses_ipv4      = ["10.0.0.0/1"]
    destination_addresses_ipv4 = ["10.0.0.0", "10.10.0.0"]

    service {
      type  = "TCP"
      value = "8081"
    }
    service {
      type  = "UDP"
      value = "22"
    }
    service {
      type  = "ICMP"
      value = "9"
    }

    description = "Bulk Rule 2 from terraform"
  }
}
