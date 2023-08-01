data "scp_tag_resources" "my_tag_resources" {
  tag_filters {
    tag_key = "okk.."
  }
}

output "contents" {
  value = data.scp_tag_resources.my_tag_resources.contents
}
