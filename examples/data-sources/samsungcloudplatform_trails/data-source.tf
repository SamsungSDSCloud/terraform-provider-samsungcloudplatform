data "samsungcloudplatform_trails" "my_trails" {

}

output "contents" {
  value = data.samsungcloudplatform_trails.my_trails.contents
}
