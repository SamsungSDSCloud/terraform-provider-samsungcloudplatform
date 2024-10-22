data "samsungcloudplatform_vpcs" "vpcs" {
}

data "samsungcloudplatform_subnets" "subnets"{
  vpc_id = data.samsungcloudplatform_vpcs.vpcs.contents[0].vpc_id
}

resource "samsungcloudplatform_vpc_dns" "vpcdns01" {
  vpc_id = data.samsungcloudplatform_vpcs.vpcs.contents[0].vpc_id
  subnet_id = data.samsungcloudplatform_subnets.subnets.contents[0].subnet_id
  domain = "hello.com"
  dns_ip = "10.254.10.254"
}
