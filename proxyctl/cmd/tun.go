package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/alexwoo79/go_coding/proxyctl/internal/proxy"
	"github.com/alexwoo79/go_coding/proxyctl/internal/tun"
	"github.com/spf13/cobra"
)

var (
	tunOnAddressFlag string
	tunOnHostFlag    string
	tunOnPortFlag    string
)

var tunCmd = &cobra.Command{
	Use:   "tun",
	Short: "管理系统级 TUN 代理（mihomo）",
	Long: `开启/关闭/查看系统级 TUN 代理（mihomo）。TUN 在网络层接管所有
应用的流量，因此无需再为 npm/pip/cargo 等工具单独配置代理。

Codex/OpenAI 与 DeepSeek 域名固定 DIRECT，不会进入 TUN。

用法示例：
  proxyctl tun on --address 10.0.0.5:7890
  proxyctl tun status
  proxyctl tun off`,
	Args: usageArgs(cobra.NoArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var tunOnCmd = &cobra.Command{
	Use:   "on",
	Short: "生成 mihomo TUN 配置并启动",
	Long: `生成受管 mihomo TUN 配置（含 OpenAI/DeepSeek DIRECT 保护规则），
用 mihomo -t 校验后通过 systemd 用户服务 mihomo-tun.service 启动。

上游地址缺省按以下顺序解析：--address/--host/--port、
PROXY_HOST/PROXY_PORT、当前系统代理、已有 mihomo 配置。`,
	Args: usageArgs(cobra.NoArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		host, port, err := resolveTUNUpstream()
		if err != nil {
			return err
		}
		if err := tun.On(host, port); err != nil {
			return fmt.Errorf("开启 TUN 失败: %w", err)
		}
		out := cmd.OutOrStdout()
		fmt.Fprintln(out, "TUN 代理已开启")
		fmt.Fprintf(out, "  上游: socks5://%s:%s\n", host, port)
		fmt.Fprintln(out, "  Codex/OpenAI 与 DeepSeek: 直连（不受 TUN 影响）")
		fmt.Fprintln(out, "  关闭请运行: proxyctl tun off")
		return nil
	},
}

var tunOffCmd = &cobra.Command{
	Use:   "off",
	Short: "停止 mihomo TUN 并清理受管配置",
	Args:  usageArgs(cobra.NoArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := tun.Off(); err != nil {
			return fmt.Errorf("关闭 TUN 失败: %w", err)
		}
		fmt.Fprintln(cmd.OutOrStdout(), "TUN 代理已关闭")
		fmt.Fprintln(cmd.OutOrStdout(), "  网络已恢复为直连（应用级代理仍可用 proxyctl on/off 控制）")
		return nil
	},
}

var tunStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "查看 mihomo TUN 状态",
	Args:  usageArgs(cobra.NoArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		st, err := tun.ReadStatus()
		if err != nil {
			return err
		}
		out := cmd.OutOrStdout()
		fmt.Fprintln(out, "=== TUN 状态 ===")
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
			return nil
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
		return nil
	},
}

// resolveTUNUpstream 解析 TUN 上游地址，缺省时依次回退到环境变量、
// 当前系统代理与已有 mihomo 配置。
func resolveTUNUpstream() (host, port string, err error) {
	host, port = strings.TrimSpace(tunOnHostFlag), strings.TrimSpace(tunOnPortFlag)
	if a := strings.TrimSpace(tunOnAddressFlag); a != "" {
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

	if host == "" || port == "" {
		if info, gerr := proxy.Get(); gerr == nil {
			switch {
			case info.HTTPEnable && info.HTTPHost != "":
				host, port = info.HTTPHost, info.HTTPPort
			case info.SOCKSEnable && info.SOCKSHost != "":
				host, port = info.SOCKSHost, info.SOCKSPort
			}
		}
	}
	if host == "" {
		if st, rerr := tun.ReadStatus(); rerr == nil && st.Upstream != "" {
			if i := strings.LastIndex(st.Upstream, ":"); i >= 0 {
				host, port = st.Upstream[:i], st.Upstream[i+1:]
			}
		}
	}

	if host == "" {
		return "", "", usageErr(errors.New("缺少上游地址：请用 --address HOST:PORT、--host HOST 或设置 PROXY_HOST"))
	}
	if port == "" {
		return "", "", usageErr(errors.New("缺少上游端口：请用 --address HOST:PORT、--port PORT 或设置 PROXY_PORT"))
	}
	if _, perr := strconv.Atoi(port); perr != nil {
		return "", "", usageErr(fmt.Errorf("无效的代理端口 %q", port))
	}
	return host, port, nil
}

func init() {
	tunCmd.AddCommand(tunOnCmd, tunOffCmd, tunStatusCmd)
	tunOnCmd.Flags().StringVarP(&tunOnAddressFlag, "address", "a", "", "上游代理地址 HOST:PORT")
	tunOnCmd.Flags().StringVarP(&tunOnHostFlag, "host", "H", "", "上游代理主机（默认 PROXY_HOST 或当前系统代理）")
	tunOnCmd.Flags().StringVarP(&tunOnPortFlag, "port", "P", "", "上游代理端口（默认 PROXY_PORT 或当前系统代理）")
}
