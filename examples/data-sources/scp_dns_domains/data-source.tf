data "scp_dns_domains" "my_scp_dns_domains" {
}

output "contents" {
  value = data.scp_dns_domains.my_scp_dns_domains.contents
}
