resource "openwrt_network_interface" "home" {
  device = "eth0"
  dns = [
    "9.9.9.9",
    "1.1.1.1",
  ]
  id      = "home"
  ipaddr  = "192.168.3.1"
  netmask = "255.255.255.0"
  proto   = "static"
}

resource "openwrt_wireless_wifi_device" "five_ghz" {
  band    = "5g"
  channel = "auto"
  id      = "radio0"
  type    = "mac80211"
}

resource "openwrt_wireless_wifi_iface" "home" {
  device                        = openwrt_wireless_wifi_device.five_ghz.id
  encryption                    = "sae"
  id                            = "wifinet0"
  key                           = "password"
  mode                          = "ap"
  network                       = openwrt_network_interface.home.id
  ssid                          = "home"
  wpa_disable_eapol_key_retries = true
}
