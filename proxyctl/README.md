# proxyctl

proxyctl 是一个系统代理与端口管理命令行工具：查看系统代理状态、一键同步 git 全局代理、执行网络连通性测试、检查端口占用并结束占用进程。

变更历史见 [CHANGELOG.md](CHANGELOG.md)。

```console
$ proxyctl status
=== 系统代理 ===
HTTP   代理: 未启用
HTTPS  代理: 未启用
SOCKS  代理: 未启用
PAC    自动代理: 未启用

=== Git 全局代理 ===
  http.proxy  = http://127.0.0.1:7892
  https.proxy = http://127.0.0.1:7892
```

## 安装

### 从源码构建

需要 Go 1.25 或更高版本：

```bash
go build -o proxyctl .
```

建议在构建时注入版本信息（版本号、commit、构建时间）。版本号可以直接取自
最新 tag：

```bash
VERSION=$(git describe --tags --match 'proxyctl/v*' --abbrev=0 2>/dev/null | sed 's|^proxyctl/||')
[ -n "$VERSION" ] || VERSION=dev
go build -ldflags "\
  -X github.com/alexwoo79/go_coding/proxyctl/cmd.version=$VERSION \
  -X github.com/alexwoo79/go_coding/proxyctl/cmd.commit=$(git rev-parse --short HEAD) \
  -X github.com/alexwoo79/go_coding/proxyctl/cmd.date=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  -o proxyctl .
```

`go install` 安装的二进制未注入版本信息，`version` 命令会显示
`dev` / `unknown`；需要带版本号的可执行文件请使用上面的 `go build` 方式。

也可以直接安装到 `$GOBIN`：

```bash
go install github.com/alexwoo79/go_coding/proxyctl@latest
```

## 发布新版本

模块位于仓库的 `proxyctl/` 子目录，因此版本 tag 必须以 `proxyctl/` 为前缀：

```bash
cd go_coding
git tag -a proxyctl/v0.1.1 -m "proxyctl v0.1.1: ..."
git push origin proxyctl/v0.1.1
```

发布后即可用 `@latest` / `@v0.1.1` 安装。如需让 `proxyctl version`
显示真实版本号，在同一 tag 上按上面的 ldflags 命令构建发布产物，例如：

```bash
VERSION=proxyctl/v0.1.1
VERSION=${VERSION#proxyctl/}
cd proxyctl
go build -ldflags "\
  -X github.com/alexwoo79/go_coding/proxyctl/cmd.version=$VERSION \
  -X github.com/alexwoo79/go_coding/proxyctl/cmd.commit=$(git rev-parse --short HEAD) \
  -X github.com/alexwoo79/go_coding/proxyctl/cmd.date=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  -o proxyctl .
./proxyctl version
```

## 快速开始

```bash
proxyctl status              # 查看系统代理与 git 代理状态
proxyctl apply               # 按系统代理自动设置 git 全局代理
proxyctl clear               # 关闭系统代理并清除 git 全局代理
proxyctl restore             # 从快照恢复 apply 前的系统代理与 git 代理
proxyctl on                  # 自动检测代理程序并一键开启
proxyctl off                 # 一键直连（等价于 clear）
proxyctl test                # 网络连通性测试
proxyctl doctor              # 一键诊断开发环境网络状态
proxyctl env                 # 生成/安装终端代理环境变量
proxyctl tools               # 管理 npm/pnpm/pip/cargo/docker/brew 的代理配置
proxyctl tun on --address 10.0.0.5:7890   # 开启系统级 TUN（mihomo）
proxyctl tun status          # 查看 TUN 状态
proxyctl tun off             # 关闭 TUN
proxyctl status --json       # 以 JSON 输出状态（便于 AI Agent 消费）
proxyctl profile list        # 列出系统代理 profile
proxyctl port 7892           # 查看 7892 端口占用
proxyctl port 7892 --kill    # 结束占用 7892 端口的进程
```

## 命令说明

| 命令 | 说明 |
| --- | --- |
| `status` | 查看系统代理（HTTP/HTTPS/SOCKS/PAC）与 git 全局代理（`http.proxy`/`https.proxy`）状态 |
| `apply` | 读取系统代理并设置 git 全局代理；优先使用 HTTP 代理，未启用时回退到 SOCKS 代理 |
| `clear` | 关闭系统代理（HTTP/HTTPS/SOCKS/PAC）并删除 git 全局代理配置 |
| `restore` | 从 `apply` 保存的状态快照恢复系统代理与 git 全局代理 |
| `on` | 自动检测运行中的代理程序（Clash/Mihomo/v2ray 等）及其端口，应用到系统代理、git、开发工具与环境变量 |
| `off` | 一键直连，等价于 `clear` |
| `test` | 依次执行公网 IP、HTTP 响应、ping、git 连通性测试 |
| `doctor` | 一键诊断：系统代理、git 代理、端口监听、连通性、环境变量，发现问题时退出码为 1 |
| `env` | 把当前系统代理生成为终端环境变量脚本，或安装到 shell 配置文件 |
| `tools` | 读写 npm/pnpm/pip/cargo/docker/brew 各自的代理配置文件（apply/clear/restore/list） |
| `tun` | 管理 mihomo 系统级 TUN 代理（on/off/status），全应用自动走代理 |
| `port` | 检查 TCP 端口占用，可结束占用进程 |
| `profile` | 保存 / 列出 / 应用 / 删除系统代理 profile |
| `version` | 显示版本、commit、构建时间与 Go 版本 |
| `completion` | 生成 bash / zsh / fish / powershell 自动补全脚本 |

运行 `proxyctl help <命令>` 可查看每个命令的详细说明。

### port

```bash
proxyctl port <端口号> [--all] [--kill] [--force]
```

| 标志 | 说明 |
| --- | --- |
| `--all` | 显示该端口的全部连接（默认仅显示 LISTEN 监听进程） |
| `--kill` | 结束占用该端口的进程 |
| `--yes` | 配合 `--kill` 跳过结束确认（非交互式环境必须使用） |
| `--force` | 配合 `--kill` 使用，强制结束进程（Unix 发送 `kill -9`，Windows 使用 `taskkill /F`） |

端口号必须是 1–65535 之间的数字。

`--kill` 默认会在终端中要求确认；在脚本等非交互式环境中必须显式追加 `--yes`。

### 状态快照与恢复

执行 `proxyctl apply` 前，proxyctl 会把当前系统代理（每个网络服务的
HTTP/HTTPS/SOCKS/PAC 状态）与 git 全局代理保存到
`~/.config/proxyctl/state.json`（可用环境变量 `PROXYCTL_STATE_FILE` 覆盖）。

```bash
proxyctl clear     # 关闭系统代理并清除 git 代理（快照保留）
proxyctl restore   # 把系统代理与 git 代理恢复为 apply 前的状态
```

这样 `clear` 不会永久丢失你原有的代理配置。

### status --json

`status --json` 输出结构化 JSON，便于脚本或 AI Agent 消费：

```json
{
  "system_proxy": {
    "http": { "enabled": true, "host": "127.0.0.1", "port": "7890" },
    "pac": { "enabled": true, "url": "http://127.0.0.1:8080/proxy.pac" }
  },
  "git": {
    "http_proxy": "http://127.0.0.1:7890",
    "https_proxy": "http://127.0.0.1:7890"
  }
}
```

未启用的代理项不会出现在输出中；未设置的 git 配置项为 `null`。

### Omarchy 上的补充状态

在 Omarchy 系统上（检测到 `omarchy` CLI 或 `/usr/share/omarchy`），
`proxyctl status` 会在常规输出后追加“Omarchy 系统设置”段落，展示当前
终端代理环境变量、会话环境文件（environment.d）、Chromium/Chrome flags，
以及 TUN/mihomo 服务、虚拟网卡、上游地址与 DIRECT 直连保护规则。

```console
$ proxyctl status
=== 系统代理 ===
HTTP   代理已启用: http://10.10.10.113:7892
...

=== Omarchy 系统设置 ===
终端代理环境变量:
  http_proxy=http://10.10.10.113:7892
  ...
会话环境文件:
  存在: /home/alex/.config/environment.d/proxy.conf
  ...
TUN/mihomo:
  mihomo-tun.service: 已停止
  虚拟网卡 Meta: 未创建
  ...
```

### doctor

`doctor` 把 `status` / `test` / `port` 组合成一次完整诊断：

```console
$ proxyctl doctor
System / Proxy / Git / Ports / Connectivity / Environment ...
Recommendation
  → 端口 7890 未监听，但 git http.proxy 指向它：请确认代理程序已启动，或运行 proxyctl clear
```

存在问题时退出码为 1，适合接入脚本与 Agent 工作流。

### profile

profile 保存的是系统代理端点配置（HTTP/HTTPS/SOCKS/PAC），存放在
`~/.config/proxyctl/profiles/`（可用环境变量 `PROXYCTL_PROFILE_DIR` 覆盖）。

```bash
proxyctl profile save clash   # 保存当前系统代理为 profile "clash"
proxyctl profile list         # 列出全部 profile
proxyctl profile use clash    # 把 profile 应用到系统代理
proxyctl profile use direct   # 关闭系统代理（直连）
proxyctl profile remove clash # 删除 profile
```

`use` 只修改系统代理；如需同步 git 代理，随后执行 `proxyctl apply`。

### env（终端环境变量代理）

macOS 系统代理只对 GUI 程序生效；终端 CLI（curl、pip、npm、cargo、uv、
brew 等）遵循环境变量 `http_proxy`/`https_proxy`/`all_proxy`/`no_proxy`。
`env` 把两者打通：

```bash
eval "$(proxyctl env)"        # 当前终端立即走代理
eval "$(proxyctl env --clear)" # 当前终端立即直连
proxyctl env install          # 写入 ~/.zshrc（可用 --file 指定），新开终端自动生效
proxyctl env remove           # 移除 install 写入的 hook
```

`install` 写入的是受管 hook，每次新开终端自动执行 `proxyctl env`，
因此会跟随当前系统代理状态。配合 profile 即可实现不同网络环境一键切换：

```bash
proxyctl profile use clash   # 切到 Clash 环境（新开终端自动走代理）
proxyctl profile use direct  # 切到直连（新开终端自动恢复直连）
```

支持 `--shell zsh|bash|sh|fish|powershell`（默认按 `$SHELL` 推断）。
首次 `install` 前会把原配置文件备份为 `<文件>.proxyctl.bak`。

### tools（开发工具代理配置）

`env` 覆盖环境变量型工具；npm、pnpm、pip、cargo、docker、brew 还有各自的
持久化配置文件，由 `tools` 管理：

```bash
proxyctl tools list                # 查看各工具当前代理配置
proxyctl tools apply               # 把当前系统代理写入全部工具（可选指定工具）
proxyctl tools clear               # 清除工具代理配置
proxyctl tools restore             # 从快照恢复（apply/clear 前自动保存快照）
proxyctl tools apply npm pnpm      # 只写指定工具
```

对应配置文件：

| 工具 | 配置文件 | 配置项 |
| --- | --- | --- |
| npm | `~/.npmrc`（可用 `npm_config_userconfig` 覆盖） | `proxy` / `https-proxy` |
| pnpm | 全局 `config.yaml`（macOS `~/Library/Preferences/pnpm/config.yaml`；Linux `~/.config/pnpm/config.yaml`；可用 `XDG_CONFIG_HOME` 覆盖） | `httpProxy` / `httpsProxy` / `noProxy` |
| pip | `~/.config/pip/pip.conf`（可用 `PIP_CONFIG_FILE` 覆盖） | `[global] proxy` |
| cargo | `~/.cargo/config.toml`（可用 `CARGO_HOME` 覆盖） | `[http] proxy` |
| docker | `~/.docker/config.json`（可用 `DOCKER_CONFIG` 覆盖） | `proxies.default.httpProxy/httpsProxy/noProxy` |
| brew | `~/.homebrew/brew.env`（可用 `XDG_CONFIG_HOME` / `HOMEBREW_XDG_CONFIG_HOME` 覆盖） | `http_proxy` / `https_proxy` / `all_proxy` / `no_proxy` |

> pnpm 的代理配置基于 pnpm v11 的全局 `config.yaml`；brew 的用户级环境文件
> `brew.env` 由 Homebrew 的 `bin/brew` 启动时读取并导出其中的代理变量，
> 无需额外安装任何东西。

`proxyctl clear` 会一并清除这些工具的代理配置；`proxyctl restore`
会一并恢复（前提是快照存在）。

### on / off 快速开关

`on` 会自动检测系统中运行中的代理程序（扫描常见端口
7890/7891/7892/1080/1087/10808/20171 等，并按进程名匹配
Clash/Mihomo/v2ray/xray/sing-box 等），然后一键应用到：

1. 系统代理（macOS 所有网络服务）
2. git 全局代理
3. 开发工具配置（npm/pnpm/pip/cargo/docker/brew）
4. 终端环境变量（新开终端自动生效，或当前终端手动 `eval "$(proxyctl env)"`）

```bash
proxyctl on    # 一键开启（执行前自动保存快照）
proxyctl off   # 一键直连（等价于 clear，快照保留可 restore）
```

未检测到代理程序时 `on` 会报错并提示先启动代理，或改用
`proxyctl profile use <名称>` 手动指定。代理程序在其他机器上时，也可直接
传入地址：

```bash
proxyctl on --address 10.10.10.113:7892
```

显式地址默认按“HTTP/HTTPS + SOCKS 同址混合端口”处理（与 proxy-on.sh 对
GNOME 系统代理的设置一致）。

若只打算用 TUN 等系统级代理、不需要 npm/pip/cargo 等工具单独走代理，
可在 `on` 时跳过开发工具配置：

```bash
proxyctl on --address 10.10.10.113:7892 --no-tools
```

### apply 与当前终端

`apply` 设置的是 git 全局配置，只对之后的 git 命令生效。子进程无法修改父 shell 的环境变量，因此若要让当前终端里的 HTTP 请求走代理，需手动执行命令输出的 `export` 提示。

## 平台支持

| 功能 | macOS | Windows | Linux + GNOME 桌面 | 其他 Unix |
| --- | --- | --- | --- | --- |
| `status` / `apply` / `clear` | ✓ | ✓（注册表） | ✓ | ✗ |
| `on` / `off` / `restore` / `profile` | ✓ | ✓ | ✓ | ✗ |
| `env` / `tools` | ✓ | ✓ | ✓ | ✓（仅文件型配置） |
| `test` | ✓ | ✓ | ✓ | ✓（无系统代理检测，回退环境变量/直连） |
| `port` | ✓ | ✓ | ✓ | ✓（需 `lsof`） |

依赖的外部命令：

- macOS：`scutil`、`networksetup`、`lsof`（系统自带）
- Windows：`netstat`、`tasklist`、`taskkill`、`powershell`（系统自带）
- Linux+GNOME：`gsettings`、`systemctl`、`dbus-update-activation-environment`
  （后两者缺失时跳过会话导入，不影响配置文件写入）、`lsof`（可用
  `apt install lsof` / `yum install lsof` 安装）
- 其他 Unix：`lsof`
- 所有平台：`git`、`ping`

`test` 命令的 HTTP 请求优先使用检测到的系统代理（HTTP → HTTPS → SOCKS），未检测到时回退到环境变量代理或直连。

### Linux 系统代理如何生效

Linux 桌面没有统一的系统代理 API，proxyctl 与 Omarchy/GNOME 环境中的
`proxy-on.sh` / `proxy-off.sh` 行为保持一致，一次性写三层：

1. `gsettings org.gnome.system.proxy`（GNOME 桌面代理）；
2. `~/.config/environment.d/proxy.conf` + `systemctl --user import-environment`
   与 `dbus-update-activation-environment`（桌面会话启动应用继承环境变量）；
3. `~/.config/chromium-flags.conf` 与 `~/.config/chrome-flags.conf`
   （Chromium/Chrome 的 `--proxy-server`，需要重启浏览器后生效）。

代理程序不在本机、而是局域网其他机器时，用显式地址而不要依赖自动检测：

```bash
proxyctl on --address 10.10.10.113:7892
proxyctl on --host 10.10.10.113 --port 7892
```

`PROXY_HOST` / `PROXY_PORT` 环境变量可作为 --host/--port 缺省值。

### Linux 辅助脚本

仓库内 [scripts/](scripts/) 提供一组薄封装（proxy-on/off/status.sh 与
proxy-tui.sh、可选 TUN 脚本），用于把 proxyctl 接入桌面快捷键或交互菜单：

```bash
cd scripts
./proxy-tui.sh status
./proxy-tui.sh on 10.10.10.113:7892
./proxy-tui.sh off
```

当前终端生效需要 `source scripts/proxy-on.sh`（或
`eval "$(proxyctl env)"`），详见 `scripts/README.md`。

### proxyctl tun（mihomo 系统级 TUN）

`proxyctl tun` 直接把 TUN 模式整合进程序：生成与 `tun-on.sh` 一致的受管
mihomo 配置（含 Codex/OpenAI、DeepSeek DIRECT 直连保护），校验后通过
systemd 用户服务 `mihomo-tun.service` 启动：

```bash
proxyctl tun on --address 10.0.0.5:7890
proxyctl tun status
proxyctl tun off
```

上游地址缺省按 `--address/--host/--port` → `PROXY_HOST/PROXY_PORT` →
当前系统代理 → 已有 mihomo 配置的顺序解析。依赖 `mihomo` 与
`systemctl --user`；如果用户单元不存在会自动创建（不会覆盖已有单元）。

开启 TUN 后所有应用由网络层接管，无需再给 npm/pip/cargo 等工具单独配置
代理；应用级代理可用 `proxyctl on --no-tools` 只设置系统层而不写工具配置。

### 在 Omarchy 上指定代理 IP/端口（样板）

```bash
# 一次性指定地址
proxyctl on --address 10.0.0.5:7890
proxyctl on --host 10.0.0.5 --port 7890

# 当前终端立即生效
eval "$(proxyctl env)"

# 查看状态
proxyctl status
```

配合薄封装脚本（source 时同时设置当前终端）：

```bash
source proxyctl/scripts/proxy-on.sh --address 10.0.0.5:7890
```

把常用地址保存为 profile 后一键切换：

```bash
proxyctl on --address 10.0.0.5:7890
proxyctl profile save home       # 保存当前端点集为 "home"
proxyctl profile use home        # 切换回 home
proxyctl profile use direct      # 直连
proxyctl apply                   # profile 切换后同步 git 代理
```

让新终端自动使用该 IP/端口（写入 `~/.bashrc` 或 `~/.zshrc`）：

```bash
export PROXY_HOST=10.0.0.5
export PROXY_PORT=7890
```

关闭：

```bash
proxyctl off
eval "$(proxyctl env --clear)"
```

只开 TUN、不写 tools 的最小流程（tun-on.sh 可直接指定上游，无需先跑
`proxyctl on`）：

```bash
./proxyctl/scripts/tun-on.sh --address 10.0.0.5:7890
./proxyctl/scripts/tun-status.sh
```

## 退出码

| 退出码 | 含义 |
| --- | --- |
| `0` | 成功 |
| `1` | 运行时错误（如无法读取系统代理、进程结束失败） |
| `2` | 用法错误（无效标志、无效参数、未知子命令） |

## 开发

```bash
go test ./...    # 运行单元测试
go vet ./...     # 静态检查
gofmt -l .       # 检查代码格式（应无输出）
```

## License

项目目前未指定开源许可证，仓库中暂无 `LICENSE` 文件。
