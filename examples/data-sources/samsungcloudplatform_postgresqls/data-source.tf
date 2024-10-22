data "samsungcloudplatform_postgresqls" "my_scp_postgresqls" {
}

output "output_my_scp_postgresqls" {
  value = data.samsungcloudplatform_postgresqls.my_scp_postgresqls
}
