package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

// 版本信息，可通过 ldflags 在构建时注入：
//
//	go build -ldflags "\
//	  -X github.com/alexwoo79/go_coding/proxyctl/cmd.version=1.0.0 \
//	  -X github.com/alexwoo79/go_coding/proxyctl/cmd.commit=$(git rev-parse --short HEAD) \
//	  -X github.com/alexwoo79/go_coding/proxyctl/cmd.date=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
//	  -o proxyctl .
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

// displayVersion 返回要展示的版本号：
// ldflags 注入的版本优先；否则读取 go install module@version 写入构建信息
// 的模块版本（如 v0.2.1）；都没有时回退 dev。
func displayVersion() string {
	if version != "" && version != "dev" {
		return version
	}
	if info, ok := debug.ReadBuildInfo(); ok {
		v := info.Main.Version
		if v != "" && v != "(devel)" {
			return v
		}
	}
	return version
}

// init 让 `proxyctl --version` / `proxyctl -v` 也使用自动解析的版本号。
func init() {
	rootCmd.Version = displayVersion()
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Args:  usageArgs(cobra.NoArgs),
	Run: func(cmd *cobra.Command, args []string) {
		out := cmd.OutOrStdout()
		fmt.Fprintf(out, "proxyctl %s (commit: %s, built: %s)\n", displayVersion(), commit, date)
		if info, ok := debug.ReadBuildInfo(); ok {
			fmt.Fprintf(out, "go: %s\n", info.GoVersion)
		}
	},
}
