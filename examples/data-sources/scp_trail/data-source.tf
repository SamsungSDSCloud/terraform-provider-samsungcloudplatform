data "scp_trails" "trails" {

}

data "scp_trail" "trail0" {
  trail_id = data.scp_trails.trails.contents[0].trail_id
}

output "trail_detail" {
  value = data.scp_trail.trail0
}
