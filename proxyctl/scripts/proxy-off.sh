#!/usr/bin/env sh
#
# proxy-off.sh — thin wrapper: disable the proxy via proxyctl.
#
#   source proxy-off.sh      # clear current terminal + system-wide
#   ./proxy-off.sh           # system-wide only

PROXY_BIN="${PROXY_BIN:-proxyctl}"

case "$0" in
    *proxy-off.sh) PROXY_SOURCED=0 ;;
    *)             PROXY_SOURCED=1 ;;
esac

exit_or_return() {
    if [ "$PROXY_SOURCED" -eq 1 ]; then
        return "${1:-0}"
    fi
    exit "${1:-0}"
}

if ! command -v "$PROXY_BIN" >/dev/null 2>&1; then
    echo "error: 未找到 $PROXY_BIN。请先构建并安装 proxyctl（go_coding/proxyctl），或设置 PROXY_BIN 指向二进制。" >&2
    exit_or_return 127
fi

"$PROXY_BIN" off || exit_or_return $?

if [ "$PROXY_SOURCED" -eq 1 ]; then
    eval "$("$PROXY_BIN" env --clear 2>/dev/null)"
    echo "  Current terminal: proxy variables cleared"
fi
