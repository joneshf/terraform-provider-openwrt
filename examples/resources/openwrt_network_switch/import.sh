# The name can be found through LuCI's web UI.
# It will be in quotes on `/cgi-bin/luci/admin/network/switch`.
# The page might say:
#     Switch "switch0"
#
# The "switch0" is the name.
# The name can also be found from LuCI's JSON-RPC API.
#
# Find the Terraform id and UCI name from LuCI's JSON-RPC API.
# One way to find this information is with `curl` and `jq`:
#
# curl \
#     --data '{"id": 0, "method": "foreach", "params": ["network", "switch"]}' \
#     http://192.168.1.1/cgi-bin/luci/rpc/uci?auth=$AUTH_TOKEN \
#     | jq '.result | map({terraformId: .[".name"], uciName: .name})'
#
# This command will output something like:
#
# [
#   {
#     "terraformId": "cfg123456",
#     "uciName": "switch0"
#   }
# ]
#
# We'd then use the information to import the appropriate resource:

terraform import openwrt_network_switch.switch0 cfg123456
