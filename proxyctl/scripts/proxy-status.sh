#!/usr/bin/env sh
#
# proxy-status.sh — thin wrapper: show proxy state via proxyctl.

PROXY_BIN="${PROXY_BIN:-proxyctl}"

if ! command -v "$PROXY_BIN" >/dev/null 2>&1; then
    echo "error: 未找到 $PROXY_BIN。请先构建并安装 proxyctl（go_coding/proxyctl），或设置 PROXY_BIN 指向二进制。" >&2
    exit 127
fi

exec "$PROXY_BIN" status
