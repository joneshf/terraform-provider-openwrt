provider "openwrt" {
  hostname = "localhost"
  port     = 8080
}

data "openwrt_network_device" "br_testing" {
  id = "br_testing"
}

output "network_device_br_testing" {
  value = data.openwrt_network_device.br_testing
}
