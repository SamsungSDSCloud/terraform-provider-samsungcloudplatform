data "samsungcloudplatform_loggingaudits" "logs" {
}

data "samsungcloudplatform_loggingaudit" "mylog" {
  logging_id = data.samsungcloudplatform_loggingaudits.logs.contents[0].id
}

output "result_scp_loggingaudits" {
  value = data.samsungcloudplatform_loggingaudit.mylog
}
