data "samsungcloudplatform_iam_policies" "my_policies" {
}

output "result_my_policies" {
  value = data.samsungcloudplatform_iam_policies.my_policies
}
