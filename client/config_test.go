package chclient

import (
	"net/http"
	"os"
	"testing"
	"time"
)

func createTempFile(t *testing.T, content string) string {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if _, err := tmpFile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	return tmpFile.Name()
}

func TestNewClientConfig(t *testing.T) {
	sampleConfig := `
fingerprint: "sample_fingerprint"
auth: "user:password"
keepalive: 30s
max-retry-count: 10
max-retry-interval: 2m
server: "example.com"
proxy: "http://proxy.com"
remotes:
  - "8080:remote->80:localhost"
headers:
  Foo: ["Bar"]
tls:
  tls-skip-verify: true
  tls-ca: "/path/to/ca"
  tls-cert: "/path/to/cert"
  tls-key: "/path/to/key"
  hostname: "example.com"
verbose: true
`

	configFile := createTempFile(t, sampleConfig)
	defer os.Remove(configFile)

	cfg, err := NewClientConfig(configFile)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	expectedHeaders := http.Header{"Foo": []string{"Bar"}}

	if cfg.Fingerprint != "sample_fingerprint" {
		t.Errorf("Expected Fingerprint to be 'sample_fingerprint', got %s", cfg.Fingerprint)
	}
	if cfg.Auth != "user:password" {
		t.Errorf("Expected Auth to be 'user:password', got %s", cfg.Auth)
	}
	if cfg.KeepAlive != 30*time.Second {
		t.Errorf("Expected KeepAlive to be 30s, got %v", cfg.KeepAlive)
	}
	if cfg.MaxRetryCount != 10 {
		t.Errorf("Expected MaxRetryCount to be 10, got %d", cfg.MaxRetryCount)
	}
	if cfg.MaxRetryInterval != 2*time.Minute {
		t.Errorf("Expected MaxRetryInterval to be 2m, got %v", cfg.MaxRetryInterval)
	}
	if cfg.Server != "example.com" {
		t.Errorf("Expected Server to be 'example.com', got %s", cfg.Server)
	}
	if cfg.Proxy != "http://proxy.com" {
		t.Errorf("Expected Proxy to be 'http://proxy.com', got %s", cfg.Proxy)
	}
	if len(cfg.Remotes) != 1 || cfg.Remotes[0] != "8080:remote->80:localhost" {
		t.Errorf("Expected Remotes to be ['8080:remote->80:localhost'], got %v", cfg.Remotes)
	}
	if len(cfg.Headers) != 1 || cfg.Headers.Get("Foo") != "Bar" {
		t.Errorf("Expected Headers to be %v, got %v", expectedHeaders, cfg.Headers)
	}
	if cfg.TLS.SkipVerify != true {
		t.Errorf("Expected TLS.SkipVerify to be true, got %v", cfg.TLS.SkipVerify)
	}
	if cfg.TLS.CA != "/path/to/ca" {
		t.Errorf("Expected TLS.CA to be '/path/to/ca', got %s", cfg.TLS.CA)
	}
	if cfg.TLS.Cert != "/path/to/cert" {
		t.Errorf("Expected TLS.Cert to be '/path/to/cert', got %s", cfg.TLS.Cert)
	}
	if cfg.TLS.Key != "/path/to/key" {
		t.Errorf("Expected TLS.Key to be '/path/to/key', got %s", cfg.TLS.Key)
	}
	if cfg.TLS.ServerName != "example.com" {
		t.Errorf("Expected TLS.ServerName to be 'example.com', got %s", cfg.TLS.ServerName)
	}
	if cfg.Verbose != true {
		t.Errorf("Expected Verbose to be true, got %v", cfg.Verbose)
	}
}

func TestNewClientConfig_DefaultFields(t *testing.T) {
	sampleConfig := `
fingerprint: "sample_fingerprint"
auth: "user:password"
server: "example.com"
proxy: "http://proxy.com"
remotes:
  - "8080:remote->80:localhost"
tls:
  tls-skip-verify: true
  tls-ca: "/path/to/ca"
  tls-cert: "/path/to/cert"
  tls-key: "/path/to/key"
  hostname: "example.com"
verbose: true
`

	configFile := createTempFile(t, sampleConfig)
	defer os.Remove(configFile)

	cfg, err := NewClientConfig(configFile)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	if cfg.KeepAlive != 25*time.Second {
		t.Errorf("Expected KeepAlive to be 30s, got %v", cfg.KeepAlive)
	}
	if cfg.MaxRetryCount != -1 {
		t.Errorf("Expected MaxRetryCount to be 10, got %d", cfg.MaxRetryCount)
	}
	if cfg.MaxRetryInterval != 0*time.Minute {
		t.Errorf("Expected MaxRetryInterval to be 2m, got %v", cfg.MaxRetryInterval)
	}
}

func TestNewClientConfig_FileNotFound(t *testing.T) {
	_, err := NewClientConfig("nonexistent.yaml")
	if err == nil {
		t.Fatal("Expected error when loading nonexistent config file, got nil")
	}
}

func TestNewClientConfig_DefaultValues(t *testing.T) {
	_, err := NewClientConfig("")
	if err != nil {
		t.Fatalf("Failed to load default config: %v", err)
	}
}
