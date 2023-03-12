# Find the Terraform id and UCI name from LuCI's JSON-RPC API.
# One way to find this information is with `curl` and `jq`:
#
# curl \
#     --data '{"id": 0, "method": "foreach", "params": ["network", "device"]}' \
#     http://192.168.1.1/cgi-bin/luci/rpc/uci?auth=$AUTH_TOKEN \
#     | jq '.result | map({terraformId: .[".name"], uciName: .name})'
#
# This command will output something like:
#
# [
#   {
#     "terraformId": "cfg030f15",
#     "uciName": "foo"
#   },
#   {
#     "terraformId": "cfg040f15",
#     "uciName": "bar"
#   }
# ]
#
# We'd then use the information to import the appropriate resource:

terraform import openwrt_network_device.foo cfg030f15
