variable "vpc_name" {
  description = "approvervpc01"
}

variable "tgw_name" {
  description = "requestertgw01"
}

variable "firewall_enable" {
  default = false
}

variable "firewall_loggable" {
  default = false
}

variable "description" {
  default = "Create TGW - VPC Connection from Terraform"
}
