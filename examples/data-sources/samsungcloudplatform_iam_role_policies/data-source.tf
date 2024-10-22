data "samsungcloudplatform_iam_roles" "my_roles1" {
}

data "samsungcloudplatform_iam_role_policies" "my_role_policies" {
  role_id = data.samsungcloudplatform_iam_roles.my_roles1.contents[0].role_id
}

output "result_my_role_policies" {
  value = data.samsungcloudplatform_iam_role_policies.my_role_policies
}
