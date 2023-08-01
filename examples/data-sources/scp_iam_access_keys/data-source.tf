data "scp_iam_access_keys" "my_access_keys" {

}

output "result_my_access_keys" {
  value = data.scp_iam_access_keys.my_access_keys
}
