data "scp_kubernetes_engine" "engine" {
  kubernetes_engine_id = "HSCLUSTER-jlJ62a08tTfPlHZTmvIocg"
}

output "output_kubernetes_engine" {
  value = data.scp_kubernetes_engine.engine
}
