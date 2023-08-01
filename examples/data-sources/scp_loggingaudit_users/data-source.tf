data "scp_loggingaudit_users" "user" {

}

output "result_scp_users" {
  value = data.scp_loggingaudit_users.user
}
