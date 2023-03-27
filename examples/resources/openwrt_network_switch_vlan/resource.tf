resource "openwrt_network_switch" "testing" {
  enable_vlan = true
  id          = "testing"
  name        = "switch0"
}

resource "openwrt_network_switch_vlan" "testing" {
  device = openwrt_network_switch.testing.name
  id     = "testing"
  ports  = "0t 1t"
  vid    = "10"
  vlan   = "2"
}
