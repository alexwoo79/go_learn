//go:build linux

package tun

import (
	"reflect"
	"strings"
	"testing"
)

func TestRenderConfig(t *testing.T) {
	content := renderConfig("10.0.0.5", "7890", false)
	for _, want := range []string{
		configMarker,
		"server: 10.0.0.5",
		"port: 7890",
		"type: socks5",
		"- DOMAIN-SUFFIX,deepseek.com,DIRECT",
		"OpenAI/ChatGPT follow the upstream proxy (default)",
		"- MATCH,GLOBAL",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("renderConfig 缺少 %q:\n%s", want, content)
		}
	}
	for _, notWant := range []string{"- DOMAIN-SUFFIX,openai.com,DIRECT", "- DOMAIN-SUFFIX,chatgpt.com,DIRECT"} {
		if strings.Contains(content, notWant) {
			t.Errorf("默认配置不应包含 %q:\n%s", notWant, content)
		}
	}
	if strings.Contains(content, "%!") {
		t.Errorf("renderConfig 存在未替换的格式动词:\n%s", content)
	}
}

func TestRenderConfigOpenAIDirect(t *testing.T) {
	content := renderConfig("10.0.0.5", "7890", true)
	for _, want := range []string{
		"- DOMAIN-SUFFIX,openai.com,DIRECT",
		"- DOMAIN-SUFFIX,chatgpt.com,DIRECT",
		"- DOMAIN-SUFFIX,oaistatic.com,DIRECT",
		`    - "*.openai.com"`,
	} {
		if !strings.Contains(content, want) {
			t.Errorf("--openai-direct 配置应包含 %q:\n%s", want, content)
		}
	}
}

func TestParseConfig(t *testing.T) {
	host, port, protected := parseConfig(renderConfig("10.0.0.5", "7890", false))
	if host != "10.0.0.5" || port != "7890" {
		t.Errorf("parseConfig host/port = %q/%q", host, port)
	}
	want := []string{"deepseek.com", "pool.ntp.org"}
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
