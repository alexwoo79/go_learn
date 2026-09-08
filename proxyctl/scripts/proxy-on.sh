#!/usr/bin/env sh
#
# proxy-on.sh — thin wrapper: enable the proxy via proxyctl.
#
# proxyctl applies the proxy to:
#   - GNOME desktop proxy (gsettings)
#   - session environment (~/.config/environment.d + systemd/D-Bus)
#   - Chromium / Chrome launch flags
#   - Git global proxy
#   - dev tools (npm/pnpm/pip/cargo/docker/brew)
#
#   source proxy-on.sh                      # current shell + system-wide
#   source proxy-on.sh --address 1.2.3.4:7890
#   ./proxy-on.sh --address 1.2.3.4:7890    # system-wide only
#
# Requires the proxyctl binary on PATH.

PROXY_BIN="${PROXY_BIN:-proxyctl}"

# Whether we were sourced (env vars can be changed in this shell) or executed.
case "$0" in
    *proxy-on.sh) PROXY_SOURCED=0 ;;
    *)            PROXY_SOURCED=1 ;;
esac

exit_or_return() {
    if [ "$PROXY_SOURCED" -eq 1 ]; then
        return "${1:-0}"
    fi
    exit "${1:-0}"
}

usage() {
    echo "Usage: source proxy-on.sh [options]"
    echo "  -a, --address HOST:PORT   Proxy address (default from environment.d or 10.10.10.113:7892)"
    echo "  -H, --host HOST           Proxy host only"
    echo "  -P, --port PORT           Proxy port only"
    echo "  -h, --help                Show this help"
}

PROXY_ADDRESS=""
PROXY_HOST=""
PROXY_PORT=""

while [ "$#" -gt 0 ]; do
    case "$1" in
        -a|--address)
            [ "$#" -ge 2 ] || { echo "error: $1 needs a value" >&2; usage >&2; exit_or_return 1; }
            PROXY_ADDRESS="$2"
            shift 2
            ;;
        -H|--host)
            [ "$#" -ge 2 ] || { echo "error: $1 needs a value" >&2; usage >&2; exit_or_return 1; }
            PROXY_HOST="$2"
            shift 2
            ;;
        -P|--port)
            [ "$#" -ge 2 ] || { echo "error: $1 needs a value" >&2; usage >&2; exit_or_return 1; }
            PROXY_PORT="$2"
            shift 2
            ;;
        -h|--help)
            usage
            exit_or_return 0
            ;;
        *)
            echo "error: unknown option: $1" >&2
            usage >&2
            exit_or_return 1
            ;;
    esac
done

# Last known address from the managed session file, otherwise the default.
DEFAULT_ADDR=$(sed -n 's/^http_proxy=http:\/\///p' \
    "${XDG_CONFIG_HOME:-$HOME/.config}/environment.d/proxy.conf" 2>/dev/null | head -n1)
[ -n "$DEFAULT_ADDR" ] || DEFAULT_ADDR="10.10.10.113:7892"

HOST=""
PORT=""
if [ -n "$PROXY_ADDRESS" ]; then
    case "$PROXY_ADDRESS" in
        *:*)
            HOST="${PROXY_ADDRESS%:*}"
            PORT="${PROXY_ADDRESS##*:}"
            ;;
        *)
            HOST="$PROXY_ADDRESS"
            ;;
    esac
else
    HOST="${DEFAULT_ADDR%:*}"
    PORT="${DEFAULT_ADDR##*:}"
fi
HOST="${PROXY_HOST:-$HOST}"
PORT="${PROXY_PORT:-$PORT}"
[ -n "$PORT" ] || PORT="${DEFAULT_ADDR##*:}"

case "$PORT" in
    ''|*[!0-9]*)
        echo "error: invalid proxy port: '$PORT'" >&2
        exit_or_return 1
        ;;
esac

if ! command -v "$PROXY_BIN" >/dev/null 2>&1; then
    echo "error: 未找到 $PROXY_BIN。请先构建并安装 proxyctl（go_coding/proxyctl），或设置 PROXY_BIN 指向二进制。" >&2
    exit_or_return 127
fi

ADDR="$HOST:$PORT"
"$PROXY_BIN" on --address "$ADDR" || exit_or_return $?

if [ "$PROXY_SOURCED" -eq 1 ]; then
    eval "$("$PROXY_BIN" env 2>/dev/null)"
    echo "  Current terminal: proxy variables set"
fi
