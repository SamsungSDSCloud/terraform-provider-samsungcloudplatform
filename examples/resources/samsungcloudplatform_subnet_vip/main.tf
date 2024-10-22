data "samsungcloudplatform_vpcs" "vpcs" {
}

data "samsungcloudplatform_subnets" "subnets"{
  vpc_id = data.samsungcloudplatform_vpcs.vpcs.contents[0].vpc_id
}


data "samsungcloudplatform_subnet_available_ips" "my_scp_subnet_available_ips1" {
  subnet_id  = data.samsungcloudplatform_subnets.subnets.contents[0].subnet_id
}

resource "samsungcloudplatform_subnet_vip" "my_subnet_vip" {
  subnet_id      = data.samsungcloudplatform_subnets.subnets.contents[0].subnet_id
  subnet_ip_id   = data.samsungcloudplatform_subnet_available_ips.my_scp_subnet_available_ips1.contents[0].ip_id
  vip_description = var.description
}
