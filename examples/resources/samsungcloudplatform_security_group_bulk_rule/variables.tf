data "terraform_remote_state" "security-group" {
  backend = "local"

  config = {
    path = "../samsungcloudplatform_security_group/terraform.tfstate"
  }
}

data "terraform_remote_state" "vpc" {
  backend = "local"

  config = {
    path = "../samsungcloudplatform_vpc/terraform.tfstate"
  }
}

variable "name" {
  default = "bulkruletestpolicy"
}
