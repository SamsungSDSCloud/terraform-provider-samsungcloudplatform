data "samsungcloudplatform_iam_policy_roles" "my_policy_roles" {
  policy_id = "policy-XXXX"
}

output "result_my_policy_roles" {
  value = data.samsungcloudplatform_iam_policy_roles.my_policy_roles
}
