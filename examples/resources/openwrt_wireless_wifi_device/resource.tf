resource "openwrt_wireless_wifi_device" "five_ghz" {
  band    = "5g"
  channel = "auto"
  id      = "cfg123456"
  type    = "mac80211"
}
