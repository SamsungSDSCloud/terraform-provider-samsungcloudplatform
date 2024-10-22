data "samsungcloudplatform_iam_group" "my_group" {

}

output "result_my_group" {
  value = data.samsungcloudplatform_iam_group.my_group
}
