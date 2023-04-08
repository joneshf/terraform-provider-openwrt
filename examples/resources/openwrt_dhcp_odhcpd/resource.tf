resource "openwrt_dhcp_odhcpd" "this" {
  id           = "testing"
  leasefile    = "/tmp/leasefile"
  leasetrigger = "/tmp/leasetrigger"
  legacy       = true
  loglevel     = 6
  maindhcp     = true
}
