provider "scp" {
}

data "scp_region" "region" {
}

variable "virtual-server-ids"{
  type = list(string)
  default = []
}

variable "numb-of-virtual-server-ids"{
  type = number
  default = 0
}

variable "placement-group-name"{
  type = string
  default = "terraform-pg"
}

variable "virtual-server-type"{
  type = string
  default = "s1"
}

variable "description" {
  type = string
  default = ""
}

data "scp_virtual_servers" "target_vm" {
  virtual_server_name = "pg-target-vm-001"
}
