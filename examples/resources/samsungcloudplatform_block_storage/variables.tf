data "terraform_remote_state" "vm" {
  backend = "local"

  config = {
    path = "../samsungcloudplatform_virtual_server/terraform.tfstate"
  }
}

variable "name" {
  default = "bstest"
}

variable "size" {
  default = "10"
}
