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
  type = string
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

variable "state" {
  type = string
  default = "RUNNING"
}

variable "ipv4" {
  type = string
  default = "192.168.29.40"
}

variable "use_dns" {
  type = bool
  default = true
}

variable "hyper_threading" {
  type = string
  default = "Y"
}

variable "nat_enabled" {
  type = bool
  default = true
}

variable "public_ip_id" {
  type = string
  default = ""
}

variable "local_subnet_enabled" {
  type = bool
  default = false
}

variable "local_subnet_id" {
  type = string
  default = ""
}

variable "local_subnet_ipv4" {
  type = string
  default = ""
}
