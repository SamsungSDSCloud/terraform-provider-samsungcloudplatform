resource "scp_internet_gateway" "my_igw" {
  vpc_id      = data.terraform_remote_state.vpc.outputs.id
  description = "Internet GW generated from Terraform"
}
