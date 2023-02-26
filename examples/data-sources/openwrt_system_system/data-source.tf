provider "openwrt" {
  hostname = "localhost"
  port     = 8080
}

data "openwrt_system_system" "this" {
}

output "system_system" {
  value = data.openwrt_system_system.this
}
