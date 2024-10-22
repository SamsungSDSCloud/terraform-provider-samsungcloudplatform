data "samsungcloudplatform_security_group_log_storages" "my_sg_log_storages" {
  vpc_id = "VPC-XXXXXXXXXXXXXXXXXXXXXX"
}

output "output_my_scp_sg_log_storages" {
  value = data.samsungcloudplatform_security_group_log_storages.my_sg_log_storages
}
