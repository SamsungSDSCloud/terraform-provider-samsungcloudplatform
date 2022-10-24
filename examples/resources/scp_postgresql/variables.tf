
data "terraform_remote_state" "vpc" {
  backend = "local"

  config = {
    path = "../scp_vpc/terraform.tfstate"
  }
}

data "terraform_remote_state" "subnet" {
  backend = "local"

  config = {
    path = "../scp_subnet/terraform.tfstate"
  }
}

data "terraform_remote_state" "security-group" {
  backend = "local"

  config = {
    path = "../scp_security_group/terraform.tfstate"
  }
}

variable "id" {
  default = "yourid"
}
variable "password" {
  default = ""
}
variable "server_name" {
  default = "svrnametest"
}

variable "cpu" {
  default = 2 #(2, 4)
}

variable "memory" {
  default = 4 #(8, 16, 32)
}
