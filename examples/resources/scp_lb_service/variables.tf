data "terraform_remote_state" "load_balancer" {
  backend = "local"

  config = {
    path = "../scp_load_balancer/terraform.tfstate"
  }
}

data "terraform_remote_state" "load_balancer_profile" {
  backend = "local"

  config = {
    path = "../scp_lb_profile/terraform.tfstate"
  }
}

variable "namel4" {
  default = "lbservicel4test"
}

variable "namel7" {
  default = "lbservicel7test"
}

