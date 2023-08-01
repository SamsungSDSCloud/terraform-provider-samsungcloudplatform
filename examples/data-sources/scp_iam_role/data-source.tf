data "scp_iam_role" "my_role" {
  role_id = "ROLE-XXXXX"
}

output "result_my_role" {
  value = data.scp_iam_role.my_role
}
