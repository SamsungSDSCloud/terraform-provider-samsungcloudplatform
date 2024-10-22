data "terraform_remote_state" "vpc" {
  backend = "local"

  config = {
    path = "../samsungcloudplatform_vpc/terraform.tfstate"
  }
}

variable "name" {
  default = "sgtest"
}
