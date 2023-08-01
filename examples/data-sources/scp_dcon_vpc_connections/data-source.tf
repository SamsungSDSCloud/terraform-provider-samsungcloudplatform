data "scp_dcon_vpc_connections" "pjt_dcon_vpc_connecions" {

}

output "contents" {
  value = data.scp_dcon_vpc_connections.pjt_dcon_vpc_connecions.contents
}
