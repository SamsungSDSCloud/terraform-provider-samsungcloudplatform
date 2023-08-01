variable "name" {
  default = "terraformbucket"
}

variable "obs_bucket_access_ip_address_ranges" {
  type = list(object({
    obs_bucket_access_ip_address_range = string
    type = string
  }))
  default = [

  ]
}
