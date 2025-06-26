data "samsungcloudplatform_kubernetes_kubeconfig" "engine" {
  kubernetes_engine_id = "HSCLUSTER-YqvFkVY8tHbN7NTlXsIyCa"
  kubeconfig_type = "private"
}

output "output_scp_kubernetes_kubeconfig" {
  value = data.samsungcloudplatform_kubernetes_kubeconfig.engine
}
