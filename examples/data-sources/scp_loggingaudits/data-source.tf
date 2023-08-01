data "scp_loggingaudits" "logs" {
}

output "result_scp_loggingaudits" {
  value = data.scp_loggingaudits.logs
}
