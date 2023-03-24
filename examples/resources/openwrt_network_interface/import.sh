# Find the Terraform id is the same as the UCI name from LuCI's JSON-RPC API.
# It is also generally the lower-cased version of the interface name in LuCI's web UI.
# One way to find this information is with `curl` and `jq`:
#
# curl \
#     --data '{"id": 0, "method": "foreach", "params": ["network", "interface"]}' \
#     http://192.168.1.1/cgi-bin/luci/rpc/uci?auth=$AUTH_TOKEN \
#     | jq '.result | map(.[".name"])'
#
# This command will output something like:
#
# [
#   "loopback",
#   "wan",
#   "wan6"
# ]
#
# We'd then use the information to import the appropriate resource:

terraform import openwrt_network_interface.loopback loopback
