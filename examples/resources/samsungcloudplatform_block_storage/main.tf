resource "samsungcloudplatform_block_storage" "my_bs" {
  name            = var.name
  storage_size_gb = var.size
  encrypt_enable  = false
  product_name    = "SSD"

  virtual_server_id = data.terraform_remote_state.vm.outputs.id
}
