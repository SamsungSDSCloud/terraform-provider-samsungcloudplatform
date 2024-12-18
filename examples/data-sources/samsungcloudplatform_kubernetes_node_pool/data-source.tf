# Find all nodepool for current project
data "samsungcloudplatform_kubernetes_node_pool" "my_scp_kubernetes_node_pool" {
  kubernetes_engine_id = "HSCLUSTER-XXXXXXXXX"
  node_pool_id = "NODEPOOL-XXXXXXXXX"
}

output "result_scp_kubernetes" {
  value = data.samsungcloudplatform_kubernetes_node_pool.my_scp_kubernetes_node_pool
}
