data "scp_gslbs" "my_scp_gslbs" {
}

output "contents" {
  value = data.scp_gslbs.my_scp_gslbs.contents
}
