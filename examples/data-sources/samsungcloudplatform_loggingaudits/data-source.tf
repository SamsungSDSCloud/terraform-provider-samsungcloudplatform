data "samsungcloudplatform_loggingaudits" "logs" {
}

output "result_scp_loggingaudits" {
  value = data.samsungcloudplatform_loggingaudits.logs
}
