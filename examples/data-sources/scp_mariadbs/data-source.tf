data "scp_mariadbs" "my_scp_mariadbs" {
}

output "output_my_scp_mariadbs" {
  value = data.scp_mariadbs.my_scp_mariadbs
}
