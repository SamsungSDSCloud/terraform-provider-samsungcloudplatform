data "terraform_remote_state" "subnet" {
  backend = "local"

  config = {
    path = "../scp_subnet/terraform.tfstate"
  }
}


data "terraform_remote_state" "public_ip" {
  backend = "local"

  config = {
    path = "../scp_public_ip/terraform.tfstate"
  }
}
