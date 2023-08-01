data "scp_iam_groups" "my_groups" {
  group_name = "ViewerGroup"
}

output "result_my_groups" {
  value = data.scp_iam_groups.my_groups
}
