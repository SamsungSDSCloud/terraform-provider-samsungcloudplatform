data "samsungcloudplatform_dcon_vpc_connections" "pjt_dcon_vpc_connecions" {

}

output "contents" {
  value = data.samsungcloudplatform_dcon_vpc_connections.pjt_dcon_vpc_connecions.contents
}
