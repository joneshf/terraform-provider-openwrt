---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "openwrt_dhcp_domain Data Source - openwrt"
subcategory: ""
description: |-
  Binds a domain name to an IP address.
---

# openwrt_dhcp_domain (Data Source)

Binds a domain name to an IP address.

## Example Usage

```terraform
data "openwrt_dhcp_domain" "testing" {
  id = "testing"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Name of the section. This name is only used when interacting with UCI directly.

### Read-Only

- `ip` (String) The IP address to be used for this domain.
- `name` (String) Hostname to assign.

