data "scp_security_group_log_storage" "my_sg_log_storage" {
  vpc_id = "VPC-XXXXXXXXXXXXXXXXXXXXXX"
  filter {
    name   = "log_storage_type"
    values = ["SECURITY_GROUP"]
  }
}

output "output_my_scp_sg_log_storage" {
  value = data.scp_security_group_log_storage.my_sg_log_storage
}
