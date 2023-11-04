data "scp_resource_group_resource_srn" "my_resource_group_resource_srn" {
  resource_id = "RESOURCE_ID-XXXXXXXX"
}

output "result_my_resource_group_resource_srn" {
  value = data.scp_resource_group_resource_srn.my_resource_group_resource_srn
}
