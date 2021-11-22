resource "docker_network" "default_static_network" {
  name   = "default_static_network"
  driver = "overlay"
}


variable "nginx_source_volume" {
  type    = string
  default = ""
}
