# Provider setup
provider "samsungcloudplatform" {
}

variable "key-pair-name" {
  type = string
  default = "terraform-keypair"
}
