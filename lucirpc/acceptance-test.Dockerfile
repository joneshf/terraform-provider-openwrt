FROM openwrtorg/rootfs:x86_64-22.03.3@sha256:bf650d3c71a5d31c51c50228c2991c6f41ef672080f911f28ce61e6ea4d54637

RUN mkdir /var/lock
RUN opkg update && opkg install \
    # Install LuCI JSON-RPC packages.
    # See https://github.com/openwrt/luci/wiki/JsonRpcHowTo#basics
    luci-compat \
    luci-lib-ipkg \
    luci-mod-rpc \
    # Install LuCI (and HTTPS support)
    # This is entirely for debugging/diagnosis purposes.
    luci \
    luci-ssl
RUN /etc/init.d/uhttpd restart
