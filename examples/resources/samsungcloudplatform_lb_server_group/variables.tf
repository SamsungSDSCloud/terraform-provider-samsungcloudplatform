data "terraform_remote_state" "load_balancer" {
  backend = "local"

  config = {
    path = "../samsungcloudplatform_load_balancer/terraform.tfstate"
  }
}

data "terraform_remote_state" "virtual_server" {
  backend = "local"

  config = {
    path = "../samsungcloudplatform_virtual_server/terraform.tfstate"
  }
}

variable "name" {
  default = "lbgrouptest"
}
