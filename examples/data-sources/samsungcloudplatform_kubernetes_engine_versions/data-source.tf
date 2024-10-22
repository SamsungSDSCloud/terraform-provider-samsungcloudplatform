# Find all available engine versions
data "samsungcloudplatform_kubernetes_engine_versions" "my_scp_kubernetes_engine_versions" {
}

output "result_scp_kubernetes_engine_versions" {
  value = data.samsungcloudplatform_kubernetes_engine_versions.my_scp_kubernetes_engine_versions
}
