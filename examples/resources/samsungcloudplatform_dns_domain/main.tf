resource "samsungcloudplatform_dns_domain" "my_dns_domain" {
  dns_domain_name       = var.name
  dns_root_domain_name  = var.root_domain_name
  dns_description       = "terraform test 2"
}
