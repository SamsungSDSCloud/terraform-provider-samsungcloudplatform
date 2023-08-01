data "scp_iam_member_groups" "my_member_groups" {
  member_id = "XXXX"
}

output "result_my_member_groups" {
  value = data.scp_iam_member_groups.my_member_groups
}
