data "terraform_remote_state" "asg" {
  backend = "local"

  config = {
    path = "../scp_auto_scaling_group/terraform.tfstate"
  }
}
variable "comparison_operator" {
  default = "GreaterThanOrEqualTo"
}
variable "comparison_operator_2" {
  default = "LessThanOrEqualTo"
}
variable "cooldown_seconds" {
  default = 60
}
variable "evaluation_minutes" {
  default = 2
}
variable "metric_method" {
  default = "AVG"
}
variable "metric_type" {
  default = "CPU"
}
variable "name" {
  default = "my-asg-policy1"
}
variable "name_2" {
  default = "my-asg-policy2"
}
variable "scale_method" {
  default = "AMOUNT"
}
variable "scale_type" {
  default = "SCALE_OUT"
}
variable "scale_type_2" {
  default = "SCALE_IN"
}
variable "scale_value" {
  default = 1
}
variable "threshold" {
  default = "60"
}
variable "threshold_2" {
  default = "20"
}
