resource "samsungcloudplatform_resource_group" "my_resource_group" {
  name = var.name

  target_resource_tags {
    tag_key = "tk01"
    tag_value = "tv01"
  }
  target_resource_tags {
    tag_key = "tk02"
    tag_value = "tv02"
  }
}
