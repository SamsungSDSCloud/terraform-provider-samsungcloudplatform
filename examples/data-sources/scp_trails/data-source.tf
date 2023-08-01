data "scp_trails" "my_trails" {

}

output "contents" {
  value = data.scp_trails.my_trails.contents
}
