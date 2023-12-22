data "terraform_remote_state" "auto_scaling_group" {
  backend = "local"

  config = {
    path = "../scp_auto_scaling_group/terraform.tfstate"
  }
}

variable "lb_rule_ids" {
  default = ["LB_RULE-XXXXX"]
}
