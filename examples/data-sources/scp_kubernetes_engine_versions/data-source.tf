# Find all available engine versions
data "scp_kubernetes_engine_versions" "my_scp_kubernetes_engine_versions" {
}

output "result_scp_kubernetes_engine_versions" {
  value = data.scp_kubernetes_engine_versions.my_scp_kubernetes_engine_versions
}
