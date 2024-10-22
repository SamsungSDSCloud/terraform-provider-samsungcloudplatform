data "samsungcloudplatform_iam_policy_groups" "my_policy_groups" {
  policy_id = "policy-XXXXX"
}

output "result_my_policy_groups" {
  value = data.samsungcloudplatform_iam_policy_groups.my_policy_groups
}
