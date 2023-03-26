resource "openwrt_network_switch" "testing" {
  enable_vlan = true
  id          = "testing"
  name        = "switch0"
}
