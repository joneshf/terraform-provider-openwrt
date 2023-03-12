resource "openwrt_network_globals" "this" {
  id              = "globals"
  packet_steering = false
  ula_prefix      = "fd12:3456:789a::/48"
}
