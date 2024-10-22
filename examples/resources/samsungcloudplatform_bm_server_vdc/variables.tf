variable "vdc_id" {
  default = "VDC-287GZVr3sLpT4HhnZVDN11"
}

variable "subnet_id" {
  default = "VDC_SUBNET-Qy0OZOsHsHmRPXyCU1W44n"
}

variable "image_name" {
  default = "Ubuntu 20.04 for BM"
}

variable "contract_discount" {
  default = "None"
}

variable "delete_protection" {
  default = false
}

variable "admin" {
  default = "root"
}

variable "password" {
  default = ""
}

variable "initial_script" {
  default = ""
}

variable "cpu" {
  default = 16
}

variable "memory" {
  default = 2048
}

locals {

  server_count = 1
  server_name_prefix = "bm-vdc-terraform-test"
  ip_address = ""
  use_hyper_threading = "N"
  dns_enabled = true

  servers = [
    # common
    for i in range(local.server_count) : {
      name = "${local.server_name_prefix}-${format("%03d", i+1)}"
      ip_address = local.ip_address
      use_hyper_threading = local.use_hyper_threading
      dns_enabled = local.dns_enabled
    }

    # separate
#    {
#      name                = "bm-vdc-test-005"
#      ip_address = "10.10.10.135"
#      use_hyper_threading = "N"
#      dns_enabled         = false
#    },
#    {
#      name                = "bm-vdc-test-006"
#      ip_address = "124"
#      use_hyper_threading = "N"
#      dns_enabled         = false
#    }
  ]
}
