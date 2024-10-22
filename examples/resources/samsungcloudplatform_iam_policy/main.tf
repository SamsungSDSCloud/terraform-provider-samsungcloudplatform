resource "samsungcloudplatform_iam_policy" "my_policy01" {
  policy_name = var.name
  policy_json = var.json
}
