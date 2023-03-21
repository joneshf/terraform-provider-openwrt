resource "openwrt_network_device" "br_testing" {
  id   = "br_testing"
  name = "br-testing"
  ports = [
    "eth0",
    "eth1",
    "eth2.20",
  ]
  type = "bridge"
}

resource "openwrt_network_interface" "testing" {
  device = openwrt_network_device.br_testing.name
  dns = [
    "9.9.9.9",
    "1.1.1.1",
  ]
  id      = "testing"
  ipaddr  = "192.168.3.1"
  netmask = "255.255.255.0"
  proto   = "static"
}
