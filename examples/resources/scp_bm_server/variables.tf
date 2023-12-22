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
  type = list(string)
  default = ["terrabm1", "terrabm2"]
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

variable "state" {
  type = list(string)
  default = ["RUNNING", "RUNNING"]
}

variable "ipv4" {
  type = list(string)
  default = ["222.222.0.2", "222.222.0.3"]
}

variable "use_dns" {
  type = list(bool)
  default = [false, false]
}

variable "hyper_threading" {
  type = list(string)
  default = ["N", "Y"]
}

variable "nat_enabled" {
  type = list(bool)
  default = [false, false]
}

variable "public_ip_id" {
  type = list(string)
  default = ["", ""]
}

variable "local_subnet_enabled" {
  type = list(bool)
  default = [false, false]
}

variable "local_subnet_id" {
  type = list(string)
  default = ["", ""]
}

variable "local_subnet_ipv4" {
  type = list(string)
  default = ["", ""]
}
