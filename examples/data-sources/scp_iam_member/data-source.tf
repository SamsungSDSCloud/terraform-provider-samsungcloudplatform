data "scp_iam_member" "my_member" {
  user_id = "XXXX"
}

output "result_my_member" {
  value = data.scp_iam_member.my_member
}
