# Find the Terraform id from LuCI's JSON-RPC API.
# One way to find this information is with `curl` and `jq`:
#
# curl \
#     --data '{"id": 0, "method": "foreach", "params": ["dhcp", "dhcp"]}' \
#     http://192.168.1.1/cgi-bin/luci/rpc/uci?auth=$AUTH_TOKEN \
#     | jq '.result | map({name: .[".name"]})'
#
# This command will output something like:
#
# [
#   {
#     "name": "lan",
#   },
#   {
#     "name": "guest",
#   }
# ]
#
# We'd then use the information to import the appropriate resource:

terraform import openwrt_dhcp_dhcp.lan lan
