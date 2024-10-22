data "samsungcloudplatform_mariadbs" "my_scp_mariadbs" {
}

output "output_my_scp_mariadbs" {
  value = data.samsungcloudplatform_mariadbs.my_scp_mariadbs
}
