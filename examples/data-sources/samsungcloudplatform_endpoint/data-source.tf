data "samsungcloudplatform_endpoint" "my_scp_endpoint" {
  endpoint_id = "ENDPOINT-XXXXXXXXXX"
#  endpoint_id = "ENDPOINT-iEYsip-Zr-pID7nkb0Abqp"
}

output "output_my_scp_endpoint" {
  value = data.samsungcloudplatform_endpoint.my_scp_endpoint
}
