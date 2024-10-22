data "terraform_remote_state" "dns_domain" {
  backend = "local"

  config = {
    path = "../../resources/scp_dns_domain/terraform.tfstate"
  }
}

data "samsungcloudplatform_dns_records" "my_scp_dns_records" {
  dns_domain_id = data.terraform_remote_state.dns_domain.outputs.id
}

output "contents" {
  value = data.samsungcloudplatform_dns_records.my_scp_dns_records.contents
}
