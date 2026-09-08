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
proxyctl status                # 在 Omarchy 上还会显示 TUN/mihomo 等系统详情
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
所有应用自动走代理。上游地址可直接用 `--address` 指定，或读取 proxyctl
生成的 `~/.config/environment.d/proxy.conf`；Codex/OpenAI 与 DeepSeek
域名固定 DIRECT，不会进入 TUN。

```bash
./tun-on.sh
./tun-on.sh --address 10.0.0.5:7890   # 直接指定上游，无需先跑 proxyctl on
./tun-status.sh
./tun-off.sh
```

开启 TUN 后所有应用都会由系统网络层接管，因此不需要再给 npm/pip/cargo
等工具单独配置代理；`proxyctl on` 配合 `--no-tools` 可跳过工具配置：

```bash
proxyctl on --address 10.0.0.5:7890 --no-tools
```

Chromium/Chrome 在启动时读取代理 flags，开关代理后需重启浏览器。
