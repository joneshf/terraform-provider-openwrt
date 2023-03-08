provider "openwrt" {
  hostname = "localhost"
  port     = 8080
}

resource "openwrt_network_globals" "this" {
  id              = "globals"
  packet_steering = false
  ula_prefix      = "fd12:3456:789a::/48"
}

output "network_globals" {
  value = resource.openwrt_network_globals.this
}
