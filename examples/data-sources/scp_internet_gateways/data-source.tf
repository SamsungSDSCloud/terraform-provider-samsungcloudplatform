data "scp_internet_gateways" "pjt_igws" {
}

output "contents" {
  value = data.scp_internet_gateways.pjt_igws.contents
}
