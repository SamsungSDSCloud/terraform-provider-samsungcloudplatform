data "samsungcloudplatform_iam_policy" "my_policy" {
  policy_id = "POLICY-XXXX"
}

output "result_my_member" {
  value = data.samsungcloudplatform_iam_policy.my_policy
}
