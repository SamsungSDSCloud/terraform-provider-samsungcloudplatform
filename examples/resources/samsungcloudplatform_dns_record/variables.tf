data "terraform_remote_state" "dns_domain" {
  backend = "local"

  config = {
    path = "../samsungcloudplatform_dns_domain/terraform.tfstate"
  }
}

variable "name" {
  default = "dns-terraform-test01"
}
