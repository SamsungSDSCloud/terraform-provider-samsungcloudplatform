data "scp_iam_roles" "my_roles1" {
}

data "scp_iam_role_policies" "my_role_policies" {
  role_id = data.scp_iam_roles.my_roles1.contents[0].role_id
}

output "result_my_role_policies" {
  value = data.scp_iam_role_policies.my_role_policies
}
