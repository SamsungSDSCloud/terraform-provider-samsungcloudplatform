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

data "scp_subnets" "subnets" {
  vpc_id = data.terraform_remote_state.vpc.outputs.id
  subnet_types = "BM"
}

variable "id" {
  default = "root"
}

variable "password" {
  default = ""
}

variable "name" {
  default = "terrabm1"
}

variable "ext_name" {
  default = "bmstorage1"
}

variable "ext_name2" {
  default = "bmstorage2"
}

variable "cpu" {
  default = 8
}

variable "memory" {
  default = 32
}
