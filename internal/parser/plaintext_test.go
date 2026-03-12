package parser

import (
	"testing"
)

func TestParsePlainText_IPPort(t *testing.T) {
	node, err := ParsePlainText("1.2.3.4:8080", "http")
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
	if node.Tag != "1.2.3.4:8080" {
		t.Errorf("expected tag 1.2.3.4:8080, got %s", node.Tag)
	}
}

func TestParsePlainText_IPPortUserPass(t *testing.T) {
	node, err := ParsePlainText("1.2.3.4:8080:admin:secret", "http")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if node.Extra["username"] != "admin" {
		t.Errorf("expected username admin, got %v", node.Extra["username"])
	}
	if node.Extra["password"] != "secret" {
		t.Errorf("expected password secret, got %v", node.Extra["password"])
	}
}

func TestParsePlainText_PasswordWithColons(t *testing.T) {
	node, err := ParsePlainText("1.2.3.4:8080:user:a:b:c", "http")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if node.Extra["password"] != "a:b:c" {
		t.Errorf("expected password a:b:c, got %v", node.Extra["password"])
	}
}

func TestParsePlainText_Socks(t *testing.T) {
	node, err := ParsePlainText("1.2.3.4:1080:user:pass", "socks")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if node.Type != "socks" {
		t.Errorf("expected type socks, got %s", node.Type)
	}
	if node.Extra["version"] != "5" {
		t.Errorf("expected version 5, got %v", node.Extra["version"])
	}
}

func TestParsePlainText_InvalidPort(t *testing.T) {
	_, err := ParsePlainText("1.2.3.4:99999", "http")
	if err == nil {
		t.Error("expected error for invalid port")
	}
}

func TestParsePlainText_InvalidFormat(t *testing.T) {
	_, err := ParsePlainText("1.2.3.4:8080:onlyuser", "http")
	if err == nil {
		t.Error("expected error for 3-part format")
	}
}

func TestParsePlainText_EmptyLine(t *testing.T) {
	_, err := ParsePlainText("", "http")
	if err == nil {
		t.Error("expected error for empty line")
	}
}

func TestParseURLWithHint_PlainText(t *testing.T) {
	node, err := ParseURLWithHint("1.2.3.4:8080:user:pass", "http")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if node.Type != "http" {
		t.Errorf("expected type http, got %s", node.Type)
	}
}

func TestParseURLWithHint_URL(t *testing.T) {
	node, err := ParseURLWithHint("http://user:pass@1.2.3.4:8080", "socks")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should use URL protocol, not hint
	if node.Type != "http" {
		t.Errorf("expected type http, got %s", node.Type)
	}
}

func TestParsePlainText_DefaultProtocol(t *testing.T) {
	node, err := ParsePlainText("1.2.3.4:8080", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if node.Type != "http" {
		t.Errorf("expected default type http, got %s", node.Type)
	}
}
