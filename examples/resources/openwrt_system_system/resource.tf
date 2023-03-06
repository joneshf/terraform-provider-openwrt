provider "openwrt" {
  hostname = "localhost"
  port     = 8080
}

resource "openwrt_system_system" "this" {
  hostname = "OpenWrt"
  id       = "cfg01e48a"
  zonename = "America/Los Angeles"
}

output "system_system" {
  value = resource.openwrt_system_system.this
}
