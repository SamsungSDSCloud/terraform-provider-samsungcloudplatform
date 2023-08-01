data "scp_direct_connects" "pjt_dcs" {

}

output "contents" {
  value = data.scp_direct_connects.pjt_dcs.contents
}
