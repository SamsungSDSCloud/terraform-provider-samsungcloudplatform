resource "scp_kubernetes_namespace" "namespace" {
  name      = var.name
  engine_id = data.terraform_remote_state.engine.outputs.id
}
