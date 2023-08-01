variable "name" {
  default = "policy_tf"
}

variable "json" {
  default = "{\"Statement\":[{\"Description\":\"sdf\",\"Effect\":\"Allow\",\"Action\":{\"iam\":{}},\"Resource\":\"*\"}]}"
}
