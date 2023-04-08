resource "openwrt_dhcp_dnsmasq" "this" {
  domain            = "testing"
  expandhosts       = true
  id                = "testing"
  local             = "/testing/"
  rebind_localhost  = true
  rebind_protection = true
}
