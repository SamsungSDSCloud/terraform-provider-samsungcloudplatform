data "terraform_remote_state" "engine" {
  backend = "local"

  config = {
    path = "../samsungcloudplatform_kubernetes_engine/terraform.tfstate"
  }
}

variable "name" {
  default = "nmtest"
}
