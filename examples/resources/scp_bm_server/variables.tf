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

// 서버 시작(RUNNING) / 중지(STOPPED)를 위한 변수
// 생성 시에는 "RUNNING"으로 세팅해야한다.
variable "state" {
  default = "STOPPED"
}
