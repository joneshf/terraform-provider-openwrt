---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "openwrt_dhcp_odhcpd Resource - openwrt"
subcategory: ""
description: |-
  An embedded DHCP/DHCPv6/RA server & NDP relay.
---

# openwrt_dhcp_odhcpd (Resource)

An embedded DHCP/DHCPv6/RA server & NDP relay.

## Example Usage

```terraform
resource "openwrt_dhcp_odhcpd" "this" {
  id           = "testing"
  leasefile    = "/tmp/leasefile"
  leasetrigger = "/tmp/leasetrigger"
  legacy       = true
  loglevel     = 6
  maindhcp     = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Name of the section. This name is only used when interacting with UCI directly.

### Optional

- `leasefile` (String) Location of the lease/hostfile for DHCPv4 and DHCPv6.
- `leasetrigger` (String) Location of the lease trigger script.
- `legacy` (Boolean) Enable DHCPv4 if the 'dhcp' section constains a `start` option, but no `dhcpv4` option set.
- `loglevel` (Number) Syslog level priority (0-7).
- `maindhcp` (Boolean) Use odhcpd as the main DHCPv4 service.

## Import

Import is supported using the following syntax:

```shell
# Find the Terraform id from LuCI's JSON-RPC API.
# One way to find this information is with `curl` and `jq`:
#
# curl \
#     --data '{"id": 0, "method": "foreach", "params": ["dhcp", "odhcpd"]}' \
#     http://192.168.1.1/cgi-bin/luci/rpc/uci?auth=$AUTH_TOKEN \
#     | jq '.result | map({terraformId: .[".name"]})'
#
# This command will output something like:
#
# [
#   {
#     "terraformId": "cfg123456",
#   }
# ]
#
# We'd then use the information to import the appropriate resource:

terraform import openwrt_dhcp_odhcpd.this cfg123456
```