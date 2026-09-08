#!/usr/bin/env sh
#
# tun-status.sh — show TUN proxy state.

echo "== TUN 服务 =="
if systemctl --user is-active --quiet mihomo-tun.service; then
    echo "  状态: 运行中"
else
    echo "  状态: 已停止"
fi

echo "== 虚拟网卡 =="
ip -br link show 2>/dev/null | rg -i '\bMeta\b' || echo "  未创建"

echo "== 上游代理 =="
CFG_FILE="${XDG_CONFIG_HOME:-$HOME/.config}/mihomo/config.yaml"
if [ -f "$CFG_FILE" ]; then
    rg '^\s+server:|^\s+port:' "$CFG_FILE" | head -2 | tr -d ' '
else
    echo "  无配置文件"
fi

echo "== 直连保护 =="
echo "  Codex/OpenAI: openai.com / chatgpt.com / chatgpt.site / oaistatic.com / oaiusercontent.com"
echo "  DeepSeek:     deepseek.com"
echo "  (以上域名在 mihomo 规则里均为 DIRECT，不会进入 TUN)"
