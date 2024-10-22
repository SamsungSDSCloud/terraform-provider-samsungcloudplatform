data "terraform_remote_state" "security-group" {
  backend = "local"

  config = {
    path = "../samsungcloudplatform_security_group/terraform.tfstate"
  }
}
