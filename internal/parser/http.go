package parser

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/xiaobei/singbox-manager/internal/storage"
)

// HttpParser HTTP Proxy Parser
type HttpParser struct{}

// Protocol Return protocol name
func (p *HttpParser) Protocol() string {
	return "http"
}

// Parse Parse HTTP proxy URL
// Format: http://user:pass@server:port#name
// Also supports https:// prefix (implies TLS enabled)
// Query params: ?tls=1&sni=...&insecure=1
func (p *HttpParser) Parse(rawURL string) (*storage.Node, error) {
	addressPart, params, name, err := parseURLParams(rawURL)
	if err != nil {
		return nil, err
	}

	var username, password, serverPart string

	// Detect if HTTPS (TLS enabled by default)
	isTLS := false
	idx := strings.Index(rawURL, "://")
	if idx != -1 {
		protocol := strings.ToLower(rawURL[:idx])
		if protocol == "https" {
			isTLS = true
		}
	}

	// Separate authentication info and server
	atIdx := strings.LastIndex(addressPart, "@")
	if atIdx != -1 {
		authPart := addressPart[:atIdx]
		serverPart = addressPart[atIdx+1:]

		if colonIdx := strings.Index(authPart, ":"); colonIdx != -1 {
			username, _ = url.QueryUnescape(authPart[:colonIdx])
			password, _ = url.QueryUnescape(authPart[colonIdx+1:])
		} else {
			decoded := tryBase64Decode(authPart)
			if decoded != "" && strings.Contains(decoded, ":") {
				colonIdx := strings.Index(decoded, ":")
				username = decoded[:colonIdx]
				password = decoded[colonIdx+1:]
			} else {
				username, _ = url.QueryUnescape(authPart)
			}
		}
	} else {
		serverPart = addressPart
	}

	// Parse server address
	server, port, err := parseServerInfo(serverPart)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse server address: %w", err)
	}

	// Set default name
	if name == "" {
		name = fmt.Sprintf("%s:%d", server, port)
	}

	// Build Extra
	extra := map[string]interface{}{}

	if username != "" {
		extra["username"] = username
	}
	if password != "" {
		extra["password"] = password
	}

	// Handle query parameter overrides
	if u := params.Get("username"); u != "" {
		extra["username"] = u
	}
	if pw := params.Get("password"); pw != "" {
		extra["password"] = pw
	}

	// TLS configuration
	if getParamBool(params, "tls") {
		isTLS = true
	}

	if isTLS {
		tls := map[string]interface{}{
			"enabled": true,
		}
		if sni := params.Get("sni"); sni != "" {
			tls["server_name"] = sni
		} else {
			tls["server_name"] = server
		}
		if getParamBool(params, "insecure") {
			tls["insecure"] = true
		}
		extra["tls"] = tls
	}

	node := &storage.Node{
		Tag:        name,
		Type:       "http",
		Server:     server,
		ServerPort: port,
		Extra:      extra,
	}

	return node, nil
}
