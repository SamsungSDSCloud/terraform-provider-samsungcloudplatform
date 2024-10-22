data "samsungcloudplatform_iam_members" "my_members" {
}

output "result_my_members" {
  value = data.samsungcloudplatform_iam_members.my_members
}
