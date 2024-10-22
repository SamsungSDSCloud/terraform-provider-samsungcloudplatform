data "samsungcloudplatform_iam_roles" "my_roles" {
}

output "result_my_roles" {
  value = data.samsungcloudplatform_iam_roles.my_roles
}
