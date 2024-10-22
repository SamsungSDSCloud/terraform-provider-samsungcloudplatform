data "samsungcloudplatform_epass" "my_scp_epass" {
}

output "output_my_scp_epass" {
  value = data.samsungcloudplatform_epass.my_scp_epass
}
