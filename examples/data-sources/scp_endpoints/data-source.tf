data "scp_endpoints" "my_scp_endpoints" {
}

output "contents" {
  value = data.scp_endpoints.my_scp_endpoints.contents
}
