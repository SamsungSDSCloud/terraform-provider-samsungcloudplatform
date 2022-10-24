data "terraform_remote_state" "vpc" {
  backend = "local"

  config = {
    path = "../scp_vpc/terraform.tfstate"
  }
}

variable "name" {
  default = "sgtest"
}
