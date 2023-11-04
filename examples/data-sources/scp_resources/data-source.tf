data "scp_resources" "my_resources" {
}

output "out_resource" {
  value = data.scp_resources.my_resources
}
