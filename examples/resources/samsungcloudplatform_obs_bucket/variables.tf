variable "name" {
  default = "terraformbucket"
}

variable "access_control_rules" {
  type = list(object({
    rule_value = string
    rule_type = string
  }))
  default = [

  ]
}
