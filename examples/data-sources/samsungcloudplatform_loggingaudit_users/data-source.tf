data "samsungcloudplatform_loggingaudit_users" "user" {

}

output "result_scp_users" {
  value = data.samsungcloudplatform_loggingaudit_users.user
}
