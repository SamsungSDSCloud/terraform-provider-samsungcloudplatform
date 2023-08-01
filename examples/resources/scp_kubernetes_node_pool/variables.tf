
data "terraform_remote_state" "engine" {
  backend = "local"

  config = {
    path = "../scp_kubernetes_engine/terraform.tfstate"
  }
}

variable "name" {
  default = "joenginenodetest"
}

variable "availability_zone_name" {
  default = "AZ1"
}
