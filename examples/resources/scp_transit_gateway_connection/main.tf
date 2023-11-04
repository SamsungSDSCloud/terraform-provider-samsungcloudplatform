data "scp_region" "my_region" {
}

resource "scp_vpc" "tgw_conn_vpc" {
  name = var.vpc_name
  description = "Approver VPC for TGW - VPC Connection"
  region = data.scp_region.my_region.location
}

resource "scp_transit_gateway" "tgw_conn_tgw" {
  transit_gateway_name = var.tgw_name
  transit_gateway_description = "Requester TGW for TGW - VPC Connection"
  region = data.scp_region.my_region.location
  bandwidth_gbps   = 1
  uplink_enabled   = false
}


resource "scp_transit_gateway_connection" "tgw_conn" {
  requester_transit_gateway_id = scp_transit_gateway.tgw_conn_tgw.id
  approver_vpc_id = scp_vpc.tgw_conn_vpc.id
  transit_gateway_connection_description = var.description
  firewall_enable = var.firewall_enable
  firewall_loggable = var.firewall_loggable
}

resource "scp_transit_gateway_connection_approve" "tgw_conn_approve" {
  transit_gateway_connection_id = scp_transit_gateway_connection.tgw_conn.id
  transit_gateway_connection_description = var.description
}
