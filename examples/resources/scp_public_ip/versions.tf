terraform {
  required_providers {
    scp = {
      version = "0.0.1"
      source  = "samsungsds/scp"
    }
  }
  required_version = ">= 0.13"
}

# Provider setup
provider "scp" {
}
