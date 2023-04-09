resource "openwrt_dhcp_host" "testing" {
  id   = "testing"
  ip   = "192.168.1.50"
  mac  = "12:34:56:78:90:ab"
  name = "testing"
}
