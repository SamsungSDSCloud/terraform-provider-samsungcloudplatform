data "scp_iam_policy" "my_policy" {
  policy_id = "POLICY-XXXX"
}

output "result_my_member" {
  value = data.scp_iam_policy.my_policy
}
