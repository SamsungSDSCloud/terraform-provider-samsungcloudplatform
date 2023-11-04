# Find public ip list
data "scp_public_ips" "my_scp_public_ips" {
}

output "contents" {
  value = data.scp_public_ips.my_scp_public_ips.contents
}
