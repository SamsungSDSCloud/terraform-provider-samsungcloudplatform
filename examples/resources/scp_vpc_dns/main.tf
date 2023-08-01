data "scp_vpcs" "vpcs" {
}

data "scp_subnets" "subnets"{
  vpc_id = data.scp_vpcs.vpcs.contents[0].vpc_id
}

resource "scp_vpc_dns" "vpcdns01" {
  vpc_id = data.scp_vpcs.vpcs.contents[0].vpc_id
  subnet_id = data.scp_subnets.subnets.contents[0].subnet_id
  domain = "hello.com"
  dns_ip = "10.254.10.254"
}
