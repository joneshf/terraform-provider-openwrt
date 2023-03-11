provider "openwrt" {
  hostname = "localhost"
  port     = 8080
}

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

output "network_device_br_testing" {
  value = resource.openwrt_network_device.br_testing
}
