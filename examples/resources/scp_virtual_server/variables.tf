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

data "terraform_remote_state" "security_group" {
  backend = "local"

  config = {
    path = "../scp_security_group/terraform.tfstate"
  }
}

data "terraform_remote_state" "public_ip" {
  backend = "local"

  config = {
    path = "../scp_public_ip/terraform.tfstate"
  }
}

variable "id" {
  default = "root"
}
variable "password" {
  default = ""
}
variable "name" {
  default = "vmtest1"
}
variable "ext_name" {
  default = "bs-vs-1"
}

variable "cpu" {
  default = 2
}

variable "memory" {
  default = 4
}
