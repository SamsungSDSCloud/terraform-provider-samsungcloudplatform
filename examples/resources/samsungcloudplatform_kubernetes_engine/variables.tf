data "terraform_remote_state" "vpc" {
  backend = "local"

  config = {
    path = "../samsungcloudplatform_vpc/terraform.tfstate"
  }
}

data "terraform_remote_state" "subnet" {
  backend = "local"

  config = {
    path = "../samsungcloudplatform_subnet/terraform.tfstate"
  }
}

data "terraform_remote_state" "security-group" {
  backend = "local"

  config = {
    path = "../samsungcloudplatform_security_group/terraform.tfstate"
  }
}

data "terraform_remote_state" "file-storage" {
  backend = "local"

  config = {
    path = "../samsungcloudplatform_file_storage/terraform.tfstate"
  }
}

data "terraform_remote_state" "load_balancer" {
  backend = "local"

  config = {
    path = "../samsungcloudplatform_load_balancer/terraform.tfstate"
  }
}

variable "name" {
  default = "enginetest"
}
