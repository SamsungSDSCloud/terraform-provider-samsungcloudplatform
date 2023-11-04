# Provider setup
provider "scp" {
}

variable "key-pair-name" {
  type = string
  default = "terraform-keypair"
}
