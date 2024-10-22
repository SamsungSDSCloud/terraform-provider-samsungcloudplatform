# Find my engines for current project
data "samsungcloudplatform_kubernetes_engines" "my_scp_kubernetes_engines" {
}

output "result_scp_kubernetes_engines" {
  value = data.samsungcloudplatform_kubernetes_engines.my_scp_kubernetes_engines
}
