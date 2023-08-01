data "scp_vpcs" "vpcs" {
}

output "contents" {
  value = data.scp_vpcs.vpcs.contents
}
