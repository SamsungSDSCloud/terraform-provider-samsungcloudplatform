data "samsungcloudplatform_trails" "trails" {

}

data "samsungcloudplatform_trail" "trail0" {
  trail_id = data.samsungcloudplatform_trails.trails.contents[0].trail_id
}

output "trail_detail" {
  value = data.samsungcloudplatform_trail.trail0
}
