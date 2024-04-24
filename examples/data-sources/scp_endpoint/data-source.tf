data "scp_endpoint" "my_scp_endpoint" {
  endpoint_id = "ENDPOINT-XXXXXXXXXX"
}

output "output_my_scp_endpoint" {
  value = data.scp_endpoint.my_scp_endpoint
}
