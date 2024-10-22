variable "name" {
  default = "mybucket"
}

data "terraform_remote_state" "vpc" {
  backend = "local"

  config = {
    path = "../samsungcloudplatform_vpc/terraform.tfstate"
  }
}
