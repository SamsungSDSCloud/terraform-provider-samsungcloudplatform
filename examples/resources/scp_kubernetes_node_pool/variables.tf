
data "terraform_remote_state" "engine" {
  backend = "local"

  config = {
    path = "../scp_kubernetes_engine/terraform.tfstate"
  }
}

variable "name" {
  default = "enginenodetest"
}
