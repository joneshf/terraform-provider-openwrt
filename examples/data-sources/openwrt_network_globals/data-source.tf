provider "openwrt" {
  hostname = "localhost"
  port     = 8080
}

data "openwrt_network_globals" "this" {
  id = "globals"
}

output "network_globals" {
  value = data.openwrt_network_globals.this
}
