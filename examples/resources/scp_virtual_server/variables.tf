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

data "terraform_remote_state" "key_pair" {
  backend = "local"

  config = {
    path = "../scp_key_pair/terraform.tfstate"
  }
}

data "terraform_remote_state" "placement_group" {
  backend = "local"

  config = {
    path = "../scp_placement_group/terraform.tfstate"
  }
}

variable "id" {
  default = "root"
}
variable "password" {
  default = ""
}
variable "name" {
  default = "vmtest2"
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
