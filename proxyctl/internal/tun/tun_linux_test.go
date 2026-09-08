//go:build linux

package tun

import (
	"reflect"
	"strings"
	"testing"
)

func TestRenderConfig(t *testing.T) {
	content := renderConfig("10.0.0.5", "7890")
	for _, want := range []string{
		configMarker,
		"server: 10.0.0.5",
		"port: 7890",
		"type: socks5",
		"- DOMAIN-SUFFIX,openai.com,DIRECT",
		"- DOMAIN-SUFFIX,chatgpt.com,DIRECT",
		"- DOMAIN-SUFFIX,deepseek.com,DIRECT",
		"- MATCH,GLOBAL",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("renderConfig 缺少 %q:\n%s", want, content)
		}
	}
	if strings.Contains(content, "%!") {
		t.Errorf("renderConfig 存在未替换的格式动词:\n%s", content)
	}
}

func TestParseConfig(t *testing.T) {
	host, port, protected := parseConfig(renderConfig("10.0.0.5", "7890"))
	if host != "10.0.0.5" || port != "7890" {
		t.Errorf("parseConfig host/port = %q/%q", host, port)
	}
	want := []string{"openai.com", "chatgpt.com", "chatgpt.site", "chatgpt-team.site", "oaistatic.com", "oaiusercontent.com", "deepseek.com", "pool.ntp.org"}
	if !reflect.DeepEqual(protected, want) {
		t.Errorf("protected = %#v, want %#v", protected, want)
	}
}

func TestParseConfigEmpty(t *testing.T) {
	host, port, protected := parseConfig("# no rules\n")
	if host != "" || port != "" || len(protected) != 0 {
		t.Errorf("空配置解析 = %q/%q/%#v", host, port, protected)
	}
}
