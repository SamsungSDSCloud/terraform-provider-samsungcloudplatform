
data "terraform_remote_state" "engine" {
  backend = "local"

  config = {
    path = "../scp_kubernetes_engine/terraform.tfstate"
  }
}

data "terraform_remote_state" "namespace" {
  backend = "local"

  config = {
    path = "../scp_kubernetes_namespace/terraform.tfstate"
  }
}

variable "name" {
  default = "appstest"
}
