data "terraform_remote_state" "security-group" {
  backend = "local"

  config = {
    path = "../scp_security_group/terraform.tfstate"
  }
}
