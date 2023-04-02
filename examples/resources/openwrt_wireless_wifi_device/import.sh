# Find the Terraform id from LuCI's JSON-RPC API.
# One way to find this information is with `curl` and `jq`:
#
# curl \
#     --data '{"id": 0, "method": "foreach", "params": ["wireless", "wifi-device"]}' \
#     http://192.168.1.1/cgi-bin/luci/rpc/uci?auth=$AUTH_TOKEN \
#     | jq '.result | map({terraformId: .[".name"]})'
#
# This command will output something like:
#
# [
#   {
#     "terraformId": "cfg123456",
#   },
#   {
#     "terraformId": "cfg123457",
#   }
# ]
#
# We'd then use the information to import the appropriate resource:

terraform import openwrt_wireless_wifi_device.five_ghz cfg123456
