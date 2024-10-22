data "samsungcloudplatform_region" "my_region" {
}

resource "samsungcloudplatform_transit_gateway_connection_approve" "tgw_conn_approve" {
  transit_gateway_connection_id = var.transit_gateway_connection_id
  transit_gateway_connection_description = var.description
}
