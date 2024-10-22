data "samsungcloudplatform_iam_access_keys" "my_access_keys" {

}

output "result_my_access_keys" {
  value = data.samsungcloudplatform_iam_access_keys.my_access_keys
}
