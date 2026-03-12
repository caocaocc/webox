package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/xiaobei/singbox-manager/internal/storage"
)

// ParsePlainText parses plain-text proxy format: IP:PORT or IP:PORT:USER:PASS
// protocolHint determines the node type: "http" (default) or "socks"
// IPv6 is not supported in plain-text format (use URL format instead).
func ParsePlainText(line string, protocolHint string) (*storage.Node, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, fmt.Errorf("empty line")
	}

	parts := strings.SplitN(line, ":", 4)

	var ip, portStr, username, password string

	switch len(parts) {
	case 2:
		// IP:PORT
		ip = parts[0]
		portStr = parts[1]
	case 4:
		// IP:PORT:USER:PASS
		ip = parts[0]
		portStr = parts[1]
		username = parts[2]
		password = parts[3]
	default:
		return nil, fmt.Errorf("invalid plain-text format, expected IP:PORT or IP:PORT:USER:PASS")
	}

	ip = strings.TrimSpace(ip)
	portStr = strings.TrimSpace(portStr)

	if ip == "" {
		return nil, fmt.Errorf("empty IP address")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		return nil, fmt.Errorf("invalid port: %s", portStr)
	}

	if protocolHint == "" {
		protocolHint = "http"
	}

	extra := map[string]interface{}{}

	nodeType := protocolHint
	if protocolHint == "socks" {
		extra["version"] = "5"
	}

	if username != "" {
		extra["username"] = username
	}
	if password != "" {
		extra["password"] = password
	}

	tag := fmt.Sprintf("%s:%d", ip, port)

	node := &storage.Node{
		Tag:        tag,
		Type:       nodeType,
		Server:     ip,
		ServerPort: port,
		Extra:      extra,
	}

	return node, nil
}
