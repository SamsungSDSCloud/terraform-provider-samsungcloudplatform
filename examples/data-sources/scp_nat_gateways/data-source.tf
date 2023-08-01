data "scp_nat_gateways" "pjt_natgws" {

}

output "contents" {
  value = data.scp_nat_gateways.pjt_natgws.contents
}
