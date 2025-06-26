data "samsungcloudplatform_kubernetes_engine" "engine" {
  kubernetes_engine_id = "HSCLUSTER-XXXXXXXXX"
}

output "output_kubernetes_engine" {
  value = data.samsungcloudplatform_kubernetes_engine.engine
}
