#!/usr/bin/env bash
#
# proxy-tui.sh — Omarchy 代理 + TUN 控制 TUI
# 菜单调用: proxy-status.sh / proxy-on.sh / proxy-off.sh /
#          tun-status.sh / tun-on.sh / tun-off.sh
#
# 用法:
#   ./proxy-tui.sh          进入交互式菜单
#   ./proxy-tui.sh status   直接显示状态
#   ./proxy-tui.sh on HOST:PORT   直接开启代理
#   ./proxy-tui.sh off      直接关闭代理

set -u

PROXY_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
STATUS_SCRIPT="$PROXY_DIR/proxy-status.sh"
ON_SCRIPT="$PROXY_DIR/proxy-on.sh"
OFF_SCRIPT="$PROXY_DIR/proxy-off.sh"
TUN_STATUS_SCRIPT="$PROXY_DIR/tun-status.sh"
TUN_ON_SCRIPT="$PROXY_DIR/tun-on.sh"
TUN_OFF_SCRIPT="$PROXY_DIR/tun-off.sh"
ENVFILE="${XDG_CONFIG_HOME:-$HOME/.config}/environment.d/proxy.conf"

for f in "$STATUS_SCRIPT" "$ON_SCRIPT" "$OFF_SCRIPT" \
         "$TUN_STATUS_SCRIPT" "$TUN_ON_SCRIPT" "$TUN_OFF_SCRIPT"; do
    [ -f "$f" ] || { echo "缺少脚本: $f" >&2; exit 1; }
done

current_address() {
    local line url
    line=$(sed -n 's/^http_proxy=http:\/\///p' "$ENVFILE" 2>/dev/null | head -n1)
    if [ -n "$line" ]; then
        url="$line"
    else
        # 脚本内置默认值
        url="10.10.10.113:7892"
    fi
    printf '%s' "$url"
}

pause_after() {
    if command -v gum >/dev/null 2>&1; then
        gum confirm --prompt.foreground=240 "按回车返回菜单" --affirmative="继续" --negative="退出" 2>/dev/null || exit 0
    else
        printf '\n按回车返回菜单...'
        read -r _ || exit 0
    fi
}

do_status() {
    bash "$STATUS_SCRIPT"
}

do_on() {
    local addr="${1:-}" default
    if [ -z "$addr" ]; then
        default=$(current_address)
        if command -v gum >/dev/null 2>&1; then
            addr=$(gum input --header "代理地址 (HOST:PORT)，留空用默认:" --placeholder "$default" --value "$default" 2>/dev/null) || addr=""
        else
            printf '代理地址 (HOST:PORT)，直接回车用默认 [%s]: ' "$default"
            read -r addr || addr=""
        fi
        [ -n "$addr" ] || addr="$default"
    fi

    case "$addr" in
        *:*)
            host=${addr%:*}
            port=${addr##*:}
            case "$port" in
                ''|*[!0-9]*)
                    echo "地址格式不对，应为 HOST:PORT，例如 10.10.10.113:7892" >&2
                    return 1
                    ;;
            esac
            ;;
        *)
            echo "地址格式不对，应为 HOST:PORT，例如 10.10.10.113:7892" >&2
            return 1
            ;;
    esac

    echo "正在开启代理: $addr"
    bash "$ON_SCRIPT" --address "$addr"
    echo
    echo "提示: 要让当前终端也立即生效，请手动执行:"
    echo "  source $ON_SCRIPT --address $addr"
}

do_off() {
    bash "$OFF_SCRIPT"
}

do_tun_status() {
    bash "$TUN_STATUS_SCRIPT"
}

do_tun_on() {
    bash "$TUN_ON_SCRIPT"
}

do_tun_off() {
    bash "$TUN_OFF_SCRIPT"
}

menu_gum() {
    local choice
    choice=$(gum choose --header "当前代理: $(current_address)" --cursor "› " \
        "1) 查看状态" \
        "2) 开启代理" \
        "3) 关闭代理" \
        "4) TUN 状态" \
        "5) TUN 开启" \
        "6) TUN 关闭" \
        "7) 退出" 2>/dev/null) || return 1
    case "$choice" in
        1*) do_status ;;
        2*) do_on ;;
        3*) do_off ;;
        4*) do_tun_status ;;
        5*) do_tun_on ;;
        6*) do_tun_off ;;
        7*) exit 0 ;;
        *) return 1 ;;
    esac
}

menu_whiptail() {
    local choice
    choice=$(whiptail --title "Omarchy 代理控制" \
        --menu "当前代理: $(current_address)" 17 60 7 \
        "1" "查看状态" \
        "2" "开启代理" \
        "3" "关闭代理" \
        "4" "TUN 状态" \
        "5" "TUN 开启" \
        "6" "TUN 关闭" \
        "7" "退出" 3>&1 1>&2 2>&3) || return 1
    case "$choice" in
        1) do_status ;;
        2) do_on ;;
        3) do_off ;;
        4) do_tun_status ;;
        5) do_tun_on ;;
        6) do_tun_off ;;
        7) exit 0 ;;
    esac
}

menu_simple() {
    printf '\n  Omarchy 代理控制\n'
    printf '  ──────────────────────────\n'
    printf '  当前代理: %s\n' "$(current_address)"
    printf '  1) 查看状态\n'
    printf '  2) 开启代理\n'
    printf '  3) 关闭代理\n'
    printf '  4) TUN 状态\n'
    printf '  5) TUN 开启\n'
    printf '  6) TUN 关闭\n'
    printf '  7) 退出\n'
    printf '  ──────────────────────────\n'
    printf '请选择 [1-7]: '
    local choice
    read -r choice || return 1
    case "$choice" in
        1) do_status ;;
        2) do_on ;;
        3) do_off ;;
        4) do_tun_status ;;
        5) do_tun_on ;;
        6) do_tun_off ;;
        7) exit 0 ;;
        *) echo "无效选择" ;;
    esac
}

# ---- 非交互模式: proxy-tui.sh status|on HOST:PORT|off ----
case "${1:-}" in
    status)      do_status; exit $? ;;
    on)          shift; do_on "${1:-}"; exit $? ;;
    off)         do_off; exit $? ;;
    tunstatus|tun-status) do_tun_status; exit $? ;;
    tunon|tun-on) do_tun_on; exit $? ;;
    tunoff|tun-off) do_tun_off; exit $? ;;
    -h|--help|help)
        sed -n '2,10p' "$0"
        exit 0
        ;;
esac

while true; do
    clear
    printf '\n  Omarchy 代理控制 TUI\n'
    if command -v gum >/dev/null 2>&1; then
        menu_gum || exit 1
    elif command -v whiptail >/dev/null 2>&1; then
        menu_whiptail || exit 1
    else
        menu_simple || exit 1
    fi
    echo
    pause_after
done
