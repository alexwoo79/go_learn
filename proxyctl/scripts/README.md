# proxyctl Linux 辅助脚本

一套薄封装脚本，用于在 Linux/GNOME 桌面上通过 `proxyctl` 开关系统代理，
并提供可选的 mihomo TUN 模式。脚本需要 `proxyctl` 已构建并位于 PATH
（可用 `PROXY_BIN` 指向自定义路径）。

## 安装

```bash
cd proxyctl
go build -o proxyctl .
cp proxyctl ~/.local/bin/

# 可选：把脚本目录加入 PATH，或复制到本地目录使用
cp -r scripts ~/proxyctl-scripts
```

## 用法

```bash
./proxy-tui.sh                 # 交互菜单（status/on/off + TUN）
./proxy-status.sh              # 查看状态
./proxy-tui.sh on 10.0.0.5:7890
./proxy-tui.sh off
```

当前终端需要 `source`，因为子进程无法修改父 shell 的环境变量：

```bash
source ~/proxyctl-scripts/proxy-on.sh --address 10.0.0.5:7890
source ~/proxyctl-scripts/proxy-off.sh
```

直接执行 `./proxy-on.sh` 只修改系统级设置（GNOME、会话环境、Git、开发
工具），不改当前 shell。

## 默认地址

`proxy-on.sh` 优先读取 proxyctl 写入的
`~/.config/environment.d/proxy.conf`，其次使用内置默认地址
`10.10.10.113:7892`；也可通过 `--address HOST:PORT`、`--host`、`--port`
或环境变量 `PROXY_HOST` / `PROXY_PORT` 覆盖。

## TUN 模式（可选）

`tun-*.sh` 依赖 `mihomo`（默认 `~/.local/bin/mihomo`），开启系统级 TUN，
所有应用自动走代理。上游地址读取 proxyctl 生成的
`~/.config/environment.d/proxy.conf`，因此与 `proxyctl on --address ...`
保持一致；Codex/OpenAI 与 DeepSeek 域名固定 DIRECT，不会进入 TUN。

```bash
./tun-on.sh
./tun-status.sh
./tun-off.sh
```

Chromium/Chrome 在启动时读取代理 flags，开关代理后需重启浏览器。
