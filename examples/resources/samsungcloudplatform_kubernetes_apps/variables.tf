data "terraform_remote_state" "engine" {
  backend = "local"

  config = {
    path = "../samsungcloudplatform_kubernetes_engine/terraform.tfstate"
  }
}

data "terraform_remote_state" "namespace" {
  backend = "local"

  config = {
    path = "../samsungcloudplatform_kubernetes_namespace/terraform.tfstate"
  }
}

variable "additional_params" {
  type = map(string)
  default = {
    "replicaCount" = "1"
  }
}

variable "name" {
  default = "tf-apps-test1"
}
