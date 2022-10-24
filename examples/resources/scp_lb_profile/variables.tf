data "terraform_remote_state" "load_balancer" {
  backend = "local"

  config = {
    path = "../scp_load_balancer/terraform.tfstate"
  }
}

variable "name_persistence" {
  default = "lbprofilepersistence"
}


variable "name" {
  default = "lbprofiletest"
}

variable "name_l7" {
  default = "lbprofilel7"
}
