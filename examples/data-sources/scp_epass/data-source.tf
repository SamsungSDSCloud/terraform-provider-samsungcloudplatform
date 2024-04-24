data "scp_epass" "my_scp_epass" {
}

output "output_my_scp_epass" {
  value = data.scp_epass.my_scp_epass
}
