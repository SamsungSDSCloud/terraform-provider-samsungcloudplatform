resource "samsungcloudplatform_dns_record" "my_dns_record" {
  dns_domain_id         = data.terraform_remote_state.dns_domain.outputs.id
  dns_record_name       = var.name
  dns_record_type       = "MX"
  ttl = 300
  dns_record_mapping {
    record_destination            = "192.168.0.1"
    preference = 1
  }
  dns_record_mapping {
    record_destination            = "192.168.0.2"
    preference = 2
  }
}
