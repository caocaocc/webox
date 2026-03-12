package parser

import (
	"testing"
)

func TestHttpParser_Basic(t *testing.T) {
	node, err := ParseURL("http://user:pass@1.2.3.4:8080#myproxy")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if node.Type != "http" {
		t.Errorf("expected type http, got %s", node.Type)
	}
	if node.Server != "1.2.3.4" {
		t.Errorf("expected server 1.2.3.4, got %s", node.Server)
	}
	if node.ServerPort != 8080 {
		t.Errorf("expected port 8080, got %d", node.ServerPort)
	}
	if node.Tag != "myproxy" {
		t.Errorf("expected tag myproxy, got %s", node.Tag)
	}
	if node.Extra["username"] != "user" {
		t.Errorf("expected username user, got %v", node.Extra["username"])
	}
	if node.Extra["password"] != "pass" {
		t.Errorf("expected password pass, got %v", node.Extra["password"])
	}
}

func TestHttpParser_NoAuth(t *testing.T) {
	node, err := ParseURL("http://1.2.3.4:3128")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if node.Type != "http" {
		t.Errorf("expected type http, got %s", node.Type)
	}
	if node.Tag != "1.2.3.4:3128" {
		t.Errorf("expected tag 1.2.3.4:3128, got %s", node.Tag)
	}
	if _, ok := node.Extra["username"]; ok {
		t.Error("expected no username")
	}
}

func TestHttpParser_HTTPS(t *testing.T) {
	node, err := ParseURL("https://user:pass@proxy.example.com:443#secure")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if node.Type != "http" {
		t.Errorf("expected type http, got %s", node.Type)
	}
	tls, ok := node.Extra["tls"].(map[string]interface{})
	if !ok {
		t.Fatal("expected tls map in extra")
	}
	if tls["enabled"] != true {
		t.Error("expected tls enabled")
	}
}

func TestHttpParser_TLSQueryParam(t *testing.T) {
	node, err := ParseURL("http://user:pass@1.2.3.4:8080?tls=1&sni=example.com&insecure=1#proxy")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tls, ok := node.Extra["tls"].(map[string]interface{})
	if !ok {
		t.Fatal("expected tls map in extra")
	}
	if tls["server_name"] != "example.com" {
		t.Errorf("expected sni example.com, got %v", tls["server_name"])
	}
	if tls["insecure"] != true {
		t.Error("expected insecure true")
	}
}

func TestHttpParser_IPv6(t *testing.T) {
	node, err := ParseURL("http://user:pass@[::1]:8080#ipv6")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if node.Server != "::1" {
		t.Errorf("expected server ::1, got %s", node.Server)
	}
}

func TestSerializeHttp(t *testing.T) {
	node, err := ParseURL("http://user:pass@1.2.3.4:8080#test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	serialized, err := SerializeNode(node)
	if err != nil {
		t.Fatalf("serialize error: %v", err)
	}
	// Re-parse to verify roundtrip
	node2, err := ParseURL(serialized)
	if err != nil {
		t.Fatalf("re-parse error: %v", err)
	}
	if node2.Server != node.Server || node2.ServerPort != node.ServerPort {
		t.Errorf("roundtrip mismatch: %s:%d vs %s:%d", node.Server, node.ServerPort, node2.Server, node2.ServerPort)
	}
	if node2.Extra["username"] != node.Extra["username"] {
		t.Errorf("username mismatch: %v vs %v", node.Extra["username"], node2.Extra["username"])
	}
}
