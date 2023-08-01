data "terraform_remote_state" "bm" {
  backend = "local"

  config = {
    path = "../scp_bm_server/terraform.tfstate"
  }
}

resource "scp_bm_block_storage" "bm_bs" {
  name = var.name
  product_name = var.product_name
  storage_size_gb = var.storage_size
  encrypted = false
  snapshot_policy = false
  snapshot_capacity_rate = var.snapshot_capacity_rate

  snap_shot_schedule = {
    "day_of_week" =  var.day_of_week
    "frequency" =  var.frequency
    "hour" =  var.hour
  }

  bm_server_ids = [data.terraform_remote_state.bm.outputs.id]
}
