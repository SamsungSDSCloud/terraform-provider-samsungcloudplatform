data "scp_loggingaudits" "logs" {
}

data "scp_loggingaudit" "mylog" {
  logging_id = data.scp_loggingaudits.logs.contents[0].id
}

output "result_scp_loggingaudits" {
  value = data.scp_loggingaudit.mylog
}
