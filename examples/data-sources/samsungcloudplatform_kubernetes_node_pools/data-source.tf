# Find all nodepool for current project
data "samsungcloudplatform_kubernetes_node_pools" "my_scp_kubernetes_node_pools" {
  kubernetes_engine_id = "engine id"
}

output "result_scp_kubernetes" {
  value = data.samsungcloudplatform_kubernetes_node_pools.my_scp_kubernetes_node_pools
}
