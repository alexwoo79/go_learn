package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/alexwoo79/go_coding/proxyctl/internal/git"
	"github.com/alexwoo79/go_coding/proxyctl/internal/proxy"
	"github.com/alexwoo79/go_coding/proxyctl/internal/toolproxy"
	"github.com/spf13/cobra"
)

var (
	onAddressFlag string
	onHostFlag    string
	onPortFlag    string
	onNoToolsFlag bool
)

var onCmd = &cobra.Command{
	Use:   "on",
	Short: "自动检测代理程序并一键开启",
	Long: `自动检测系统中运行中的代理程序（Clash/Mihomo/v2ray/xray/sing-box 等）
及其端口，然后把代理应用到：
  1. 系统代理（桌面/系统级 HTTP/SOCKS 设置）
  2. git 全局代理（http.proxy / https.proxy）
  3. 开发工具配置（npm / pnpm / pip / cargo / docker / brew）
  4. 终端环境变量（新开终端自动生效，或手动 eval）
	执行前会保存快照，可用 proxyctl off / restore 恢复。

代理程序不在本机运行时（例如上游在局域网其他机器），可用显式地址：
  proxyctl on --address 10.10.10.113:7892
  proxyctl on --host 10.10.10.113 --port 7892

如果只用 TUN 等系统级代理、不希望改动 npm/pip/cargo 等工具的配置文件，
可加 --no-tools 跳过开发工具配置：
  proxyctl on --address 10.10.10.113:7892 --no-tools`,
	Args: usageArgs(cobra.NoArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			d        proxy.Detected
			err      error
			explicit bool
		)
		if onAddressFlag != "" || onHostFlag != "" || onPortFlag != "" {
			d, err = detectedFromFlags()
			explicit = err == nil
		} else {
			d, err = proxy.Detect()
		}
		if err != nil {
			return err
		}
		out := cmd.OutOrStdout()
		fmt.Fprintf(out, "检测到代理程序: %s (%s)\n", d.Command, d.URL())

		// 1. 先保存快照（系统 + git；未跳过时含 tools），保证 off / restore 可恢复
		path, err := saveStateSnapshot()
		if err != nil {
			return err
		}
		if !onNoToolsFlag {
			if err := snapshotToolConfigs(toolproxy.Supported()); err != nil {
				return err
			}
		}

		// 2. 设置系统代理
		var httpEP, httpsEP, socksEP *proxy.EndpointState
		if explicit && d.Scheme != "socks5h" {
			// 显式地址默认按混合端口处理：HTTP/HTTPS 与 SOCKS 都指向同一
			// 地址，与旧 proxy-on.sh 对 GNOME 系统代理的设置一致。
			httpEP = &proxy.EndpointState{Enabled: true, Host: d.Host, Port: d.Port}
			httpsEP = &proxy.EndpointState{Enabled: true, Host: d.Host, Port: d.Port}
			socksEP = &proxy.EndpointState{Enabled: true, Host: d.Host, Port: d.Port}
		} else if d.Scheme == "socks5h" {
			socksEP = &proxy.EndpointState{Enabled: true, Host: d.Host, Port: d.Port}
		} else {
			httpEP = &proxy.EndpointState{Enabled: true, Host: d.Host, Port: d.Port}
		}
		if err := proxy.ApplyProfile(httpEP, httpsEP, socksEP, nil); err != nil {
			return fmt.Errorf("设置系统代理失败: %w", err)
		}
		fmt.Fprintf(out, "已设置系统代理: %s\n", d.URL())

		// 3. 设置 git 代理
		url := d.URL()
		gitCfg := git.New()
		if err := gitCfg.Set("http.proxy", url); err != nil {
			return fmt.Errorf("设置 git http.proxy 失败: %w", err)
		}
		if err := gitCfg.Set("https.proxy", url); err != nil {
			return fmt.Errorf("设置 git https.proxy 失败: %w", err)
		}
		fmt.Fprintf(out, "已设置 git 代理: %s\n", url)

		// 4. 应用开发工具配置（--no-tools 时跳过）
		if !onNoToolsFlag {
			if err := toolproxy.ApplyTo(toolproxy.Supported(), url); err != nil {
				return fmt.Errorf("应用开发工具代理失败: %w", err)
			}
			fmt.Fprintf(out, "已应用开发工具代理：%s\n", strings.Join(toolNames(toolproxy.Supported()), ", "))
		} else {
			fmt.Fprintln(out, "已跳过开发工具代理配置（--no-tools）")
		}

		// 5. 提示终端环境变量
		fmt.Fprintf(out, "已保存快照: %s\n", path)
		fmt.Fprintln(out, "当前终端如需立即生效，请执行: eval \"$(proxyctl env)\"")
		fmt.Fprintln(out, "新开终端自动生效可先运行: proxyctl env install")
		return nil
	},
}

// detectedFromFlags 根据 --address/--host/--port 构造显式代理端点。
// host/port 允许通过 PROXY_HOST / PROXY_PORT 环境变量覆盖（与旧脚本一致）。
func detectedFromFlags() (proxy.Detected, error) {
	host, port := strings.TrimSpace(onHostFlag), strings.TrimSpace(onPortFlag)
	if a := strings.TrimSpace(onAddressFlag); a != "" {
		if strings.Contains(a, ":") {
			if i := strings.LastIndex(a, ":"); i >= 0 {
				host, port = strings.TrimSpace(a[:i]), strings.TrimSpace(a[i+1:])
			}
		} else {
			host = a
		}
	}
	if host == "" {
		host = strings.TrimSpace(os.Getenv("PROXY_HOST"))
	}
	if port == "" {
		port = strings.TrimSpace(os.Getenv("PROXY_PORT"))
	}
	if host == "" {
		return proxy.Detected{}, usageErr(fmt.Errorf("缺少代理主机：请用 --address HOST:PORT、--host HOST 或设置 PROXY_HOST"))
	}
	if port == "" {
		port = "7892"
	}
	if _, err := strconv.Atoi(port); err != nil {
		return proxy.Detected{}, usageErr(fmt.Errorf("无效的代理端口 %q", port))
	}
	return proxy.Detected{
		Scheme:  "http",
		Host:    host,
		Port:    port,
		Command: "手动指定",
		Source:  "命令行参数",
	}, nil
}

func init() {
	onCmd.Flags().StringVarP(&onAddressFlag, "address", "a", "", "显式代理地址 HOST:PORT（HTTP；可配合 PROXY_HOST/PROXY_PORT）")
	onCmd.Flags().StringVarP(&onHostFlag, "host", "H", "", "显式代理主机（默认取自 PROXY_HOST）")
	onCmd.Flags().StringVarP(&onPortFlag, "port", "P", "", "显式代理端口（默认取自 PROXY_PORT 或 7892）")
	onCmd.Flags().BoolVar(&onNoToolsFlag, "no-tools", false, "不写入 npm/pnpm/pip/cargo/docker/brew 的代理配置")
}

var offCmd = &cobra.Command{
	Use:   "off",
	Short: "一键直连（等价于 clear）",
	Long: `关闭系统代理（含 WPAD）、清除 git 与开发工具（npm/pnpm/pip/cargo/docker/brew）
的代理配置，并提示当前终端执行 eval "$(proxyctl env --clear)" 立即直连。`,
	Args: usageArgs(cobra.NoArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runClear(cmd)
	},
}
