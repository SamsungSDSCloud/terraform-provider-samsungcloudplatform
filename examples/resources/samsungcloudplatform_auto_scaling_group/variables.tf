data "terraform_remote_state" "vpc" {
  backend = "local"

  config = {
    path = "../samsungcloudplatform_vpc/terraform.tfstate"
  }
}

data "terraform_remote_state" "subnet" {
  backend = "local"

  config = {
    path = "../samsungcloudplatform_subnet/terraform.tfstate"
  }
}

data "terraform_remote_state" "security_group" {
  backend = "local"

  config = {
    path = "../samsungcloudplatform_security_group/terraform.tfstate"
  }
}
data "terraform_remote_state" "lc" {
  backend = "local"

  config = {
    path = "../samsungcloudplatform_launch_configuration/terraform.tfstate"
  }
}

variable "az_name" {
  default = "AZ1"
}
variable "name" {
  default = "my_asg"
}
variable "desired" {
  default = 0
}
variable "min" {
  default = 0
}
variable "max" {
  default = 0
}
