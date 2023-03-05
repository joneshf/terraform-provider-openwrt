provider "openwrt" {
  hostname = "localhost"
  port     = 8080
}

data "openwrt_system_system" "this" {
  id = "cfg01e48a"
}

output "system_system" {
  value = data.openwrt_system_system.this
}
