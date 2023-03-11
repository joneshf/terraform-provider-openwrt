---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "openwrt_network_device Data Source - openwrt"
subcategory: ""
description: |-
  A physical or virtual "device" in OpenWrt jargon. Commonly referred to as an "interface" in other networking jargon.
---

# openwrt_network_device (Data Source)

A physical or virtual "device" in OpenWrt jargon. Commonly referred to as an "interface" in other networking jargon.

## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Name of the section. This name is only used when interacting with UCI directly.

### Read-Only

- `bridge_empty` (Boolean) Bring up the bridge device even if no ports are attached
- `dadtransmits` (Number) Amount of Duplicate Address Detection probes to send
- `ipv6` (Boolean) Enable IPv6 for the device.
- `macaddr` (String) MAC Address of the device.
- `mtu` (Number) Maximum Transmissible Unit.
- `mtu6` (Number) Maximum Transmissible Unit for IPv6.
- `name` (String) Name of the device. This name is referenced in other network configuration.
- `ports` (Set of String) Specifies the wired ports to attach to this bridge.
- `txqueuelen` (Number) Transmission queue length.
- `type` (String) The type of device. Currently, only "bridge" is supported.

