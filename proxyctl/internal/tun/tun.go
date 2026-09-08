// Package tun 封装 mihomo 系统级 TUN 代理的管理：
// 生成受管配置、校验、启动/停止 systemd 用户服务与读取状态。
//
// TUN 模式只在 Linux 上实现；其他平台 On/Off/ReadStatus 返回明确错误。
package tun

// Status 是 TUN/mihomo 的只读状态。
type Status struct {
	ConfigPath   string   `json:"config_path,omitempty"`
	ConfigExists bool     `json:"config_exists"`
	Service      string   `json:"service"` // running / stopped / systemctl 不可用
	Upstream     string   `json:"upstream,omitempty"`
	Interfaces   []string `json:"interfaces,omitempty"`
	Protected    []string `json:"protected_domains,omitempty"`
}
