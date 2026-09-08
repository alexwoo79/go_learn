//go:build linux

package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alexwoo79/go_coding/proxyctl/internal/tun"
)

// proxyEnvVarNames 是 status 关心的终端代理环境变量。
var proxyEnvVarNames = []string{
	"http_proxy", "https_proxy", "all_proxy", "no_proxy",
	"HTTP_PROXY", "HTTPS_PROXY", "ALL_PROXY", "NO_PROXY",
}

// printOmarchyStatus 在 Omarchy 系统上输出与桌面/系统设置相关的补充状态：
// 终端环境变量、会话环境文件、Chromium/Chrome 启动参数与 TUN/mihomo 状态。
// 其他系统上该函数为空操作（见 status_other.go）。
func printOmarchyStatus(out io.Writer) {
	if !omarchyDetected() {
		return
	}

	fmt.Fprintln(out)
	fmt.Fprintln(out, "=== Omarchy 系统设置 ===")

	// 1. 当前终端环境变量
	fmt.Fprintln(out, "终端代理环境变量:")
	printed := false
	for _, name := range proxyEnvVarNames {
		if v := os.Getenv(name); v != "" {
			fmt.Fprintf(out, "  %s=%s\n", name, v)
			printed = true
		}
	}
	if !printed {
		fmt.Fprintln(out, "  （未设置）")
	}

	// 2. 会话环境文件（桌面会话启动的应用继承）
	envFile := linuxConfigPath("environment.d", "proxy.conf")
	fmt.Fprintln(out, "会话环境文件:")
	if data, err := os.ReadFile(envFile); err == nil {
		fmt.Fprintf(out, "  存在: %s\n", envFile)
		for _, line := range strings.Split(strings.TrimRight(string(data), "\n"), "\n") {
			if line != "" {
				fmt.Fprintf(out, "    %s\n", line)
			}
		}
	} else {
		fmt.Fprintln(out, "  不存在")
	}

	// 3. Chromium / Chrome 启动参数
	fmt.Fprintln(out, "Chromium/Chrome flags:")
	for _, name := range []string{"chromium-flags.conf", "chrome-flags.conf"} {
		path := linuxConfigPath(name)
		if data, err := os.ReadFile(path); err == nil {
			text := strings.TrimSpace(string(data))
			if text == "" {
				fmt.Fprintf(out, "  %s: （空文件）\n", path)
			} else {
				fmt.Fprintf(out, "  %s:\n", path)
				for _, line := range strings.Split(text, "\n") {
					if line != "" {
						fmt.Fprintf(out, "    %s\n", line)
					}
				}
			}
		} else {
			fmt.Fprintf(out, "  %s: 未创建\n", path)
		}
	}

	// 4. TUN / mihomo
	printTUNStatus(out)
}

// printTUNStatus 输出 mihomo TUN 服务、虚拟网卡、上游与直连保护规则。
func printTUNStatus(out io.Writer) {
	st, err := tun.ReadStatus()
	if err != nil {
		fmt.Fprintf(out, "TUN/mihomo: 读取失败: %v\n", err)
		return
	}
	fmt.Fprintln(out, "TUN/mihomo:")
	switch st.Service {
	case "running":
		fmt.Fprintln(out, "  mihomo-tun.service: 运行中")
	case "systemctl 不可用":
		fmt.Fprintf(out, "  mihomo-tun.service: %s\n", st.Service)
	default:
		fmt.Fprintln(out, "  mihomo-tun.service: 已停止")
	}
	if len(st.Interfaces) == 0 {
		fmt.Fprintln(out, "  虚拟网卡 Meta: 未创建")
	} else {
		fmt.Fprintf(out, "  虚拟网卡 Meta: %s\n", strings.Join(st.Interfaces, ", "))
	}
	if !st.ConfigExists {
		fmt.Fprintln(out, "  mihomo 配置: 不存在")
		return
	}
	fmt.Fprintf(out, "  mihomo 配置: %s\n", st.ConfigPath)
	if st.Upstream != "" {
		fmt.Fprintf(out, "  上游: %s\n", st.Upstream)
	} else {
		fmt.Fprintln(out, "  上游: 未找到 server/port")
	}
	if len(st.Protected) > 0 {
		fmt.Fprintln(out, "  直连保护（DIRECT）:")
		for _, d := range st.Protected {
			fmt.Fprintf(out, "    %s\n", d)
		}
	}
}

// omarchyDetected 判断当前系统是否为 Omarchy：
// omarchy CLI 在 PATH 上，或 /usr/share/omarchy 存在。
func omarchyDetected() bool {
	if _, err := exec.LookPath("omarchy"); err == nil {
		return true
	}
	if _, err := os.Stat("/usr/share/omarchy"); err == nil {
		return true
	}
	return false
}

// linuxConfigPath 返回 XDG 配置目录下的路径（未设置时回退 ~/.config）。
func linuxConfigPath(parts ...string) string {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			home = "."
		}
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(append([]string{base}, parts...)...)
}
