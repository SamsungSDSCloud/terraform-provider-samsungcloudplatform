data "scp_vpcs" "vpcs" {
}

data "scp_vpc_dnss" "dnss" {
  vpc_id = data.scp_vpcs.vpcs.contents[0].vpc_id
}

output "contents" {
  value = data.scp_vpc_dnss.dnss.contents
}
