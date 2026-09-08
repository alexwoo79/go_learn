// Package proxy 封装了读取与管理系统代理配置的逻辑。
//
// 各平台实现不同：
//   - macOS 通过 scutil/networksetup；
//   - Windows 通过注册表；
//   - Linux（GNOME 桌面）通过 gsettings + ~/.config/environment.d +
//     Chromium/Chrome flags；
//   - 其他 Unix 平台不支持系统代理管理。
package proxy

import "fmt"

// Info 表示系统代理配置信息。
type Info struct {
	HTTPEnable  bool
	HTTPHost    string
	HTTPPort    string
	HTTPSEnable bool
	HTTPSHost   string
	HTTPSPort   string
	SOCKSEnable bool
	SOCKSHost   string
	SOCKSPort   string
	AutoConfig  bool
	// AutoConfigURL 是 PAC 自动代理脚本地址（ProxyAutoConfigURLString）。
	AutoConfigURL string
	// AutoDiscovery 表示是否启用代理自动发现（WPAD，仅 macOS）。
	AutoDiscovery bool
}

// HTTPProxyURL 返回 HTTP 代理 URL，未启用时返回空字符串。
func (i *Info) HTTPProxyURL() string {
	if !i.HTTPEnable {
		return ""
	}
	return i.proxyURL("http", i.HTTPHost, i.HTTPPort)
}

// HTTPSProxyURL 返回 HTTPS 代理 URL，未启用时返回空字符串。
func (i *Info) HTTPSProxyURL() string {
	if !i.HTTPSEnable {
		return ""
	}
	return i.proxyURL("http", i.HTTPSHost, i.HTTPSPort)
}

// SOCKSPProxyURL 返回 SOCKS5 代理 URL，未启用时返回空字符串。
func (i *Info) SOCKSPProxyURL() string {
	if !i.SOCKSEnable {
		return ""
	}
	return i.proxyURL("socks5h", i.SOCKSHost, i.SOCKSPort)
}

// PACURL 返回 PAC 自动代理脚本地址，未启用时返回空字符串。
func (i *Info) PACURL() string {
	if !i.AutoConfig {
		return ""
	}
	return i.AutoConfigURL
}

// proxyURL 拼接代理 URL；主机为空时返回空字符串，端口为空时省略端口。
func (i *Info) proxyURL(scheme, host, port string) string {
	if host == "" {
		return ""
	}
	if port == "" {
		return scheme + "://" + host
	}
	return fmt.Sprintf("%s://%s:%s", scheme, host, port)
}

// SystemSnapshot 是 clear/restore 需要保留的系统代理状态。
// 各平台填充不同字段：macOS 用 Services，Windows 用
// ProxyEnable/ProxyServer/AutoConfigURL，Linux 用 Linux。
type SystemSnapshot struct {
	Services      []ServiceState `json:"services,omitempty"`
	ProxyEnable   bool           `json:"proxy_enable"`
	ProxyServer   string         `json:"proxy_server,omitempty"`
	AutoConfigURL string         `json:"auto_config_url,omitempty"`
	Linux         *LinuxSnapshot `json:"linux,omitempty"`
}

// LinuxSnapshot 是 Linux（GNOME 桌面）系统代理的完整状态快照：
// gsettings 桌面代理、会话环境文件以及 Chromium/Chrome 启动参数。
type LinuxSnapshot struct {
	Desktop  *LinuxDesktopSnapshot `json:"desktop,omitempty"`
	EnvFile  *LinuxEnvFileSnapshot `json:"env_file,omitempty"`
	Chromium map[string]bool       `json:"chromium_files,omitempty"`
}

// LinuxDesktopSnapshot 记录 GNOME org.gnome.system.proxy 的完整状态。
type LinuxDesktopSnapshot struct {
	Mode        string         `json:"mode"`
	HTTP        *EndpointState `json:"http,omitempty"`
	HTTPS       *EndpointState `json:"https,omitempty"`
	SOCKS       *EndpointState `json:"socks,omitempty"`
	PAC         *PACState      `json:"pac,omitempty"`
	IgnoreHosts []string       `json:"ignore_hosts,omitempty"`
}

// LinuxEnvFileSnapshot 记录会话环境文件（environment.d/proxy.conf）在
// proxyctl 修改前的状态，供 restore 精确恢复。
type LinuxEnvFileSnapshot struct {
	Existed bool   `json:"existed"`
	Content string `json:"content,omitempty"`
}

// ServiceState 记录单个网络服务的代理状态（macOS）。
type ServiceState struct {
	Name          string         `json:"name"`
	HTTP          *EndpointState `json:"http,omitempty"`
	HTTPS         *EndpointState `json:"https,omitempty"`
	SOCKS         *EndpointState `json:"socks,omitempty"`
	PAC           *PACState      `json:"pac,omitempty"`
	AutoDiscovery *bool          `json:"auto_discovery,omitempty"`
}

// EndpointState 记录单个 HTTP/HTTPS/SOCKS 代理端点的状态。
type EndpointState struct {
	Enabled bool   `json:"enabled"`
	Host    string `json:"host,omitempty"`
	Port    string `json:"port,omitempty"`
}

// PACState 记录 PAC 自动代理的状态与脚本地址。
type PACState struct {
	Enabled bool   `json:"enabled"`
	URL     string `json:"url,omitempty"`
}
