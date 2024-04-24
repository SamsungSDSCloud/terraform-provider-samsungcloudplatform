terraform {
  required_providers {
    scp = {
      version = "3.5.1"
      source = "SamsungSDSCloud/samsungcloudplatform"
    }
  }
  required_version = ">=0.13"
}

provider "scp" {

}
