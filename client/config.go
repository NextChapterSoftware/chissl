package chclient

import (
	"context"
	"errors"
	yaml "gopkg.in/yaml.v3"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"time"
)

// Config represents a client configuration
type Config struct {
	Fingerprint      string        `yaml:"fingerprint,omitempty"`
	Auth             string        `yaml:"auth,omitempty"`
	KeepAlive        time.Duration `yaml:"keepalive,omitempty"`
	MaxRetryCount    int           `yaml:"max-retry-count,omitempty"`
	MaxRetryInterval time.Duration `yaml:"max-retry-interval,omitempty"`
	Server           string        `yaml:"server,omitempty"`
	Proxy            string        `yaml:"proxy,omitempty"`
	Remotes          []string      `yaml:"remotes,omitempty"`
	Headers          http.Header   `yaml:"headers,omitempty"`
	TLS              TLSConfig     `yaml:"tls,omitempty"`
	DialContext      func(ctx context.Context, network, addr string) (net.Conn, error)
	Verbose          bool `yaml:"verbose,omitempty"`
}

// TLSConfig for a Client
type TLSConfig struct {
	SkipVerify bool   `yaml:"tls-skip-verify,omitempty"`
	CA         string `yaml:"tls-ca,omitempty"`
	Cert       string `yaml:"tls-cert,omitempty"`
	Key        string `yaml:"tls-key,omitempty"`
	ServerName string `yaml:"hostname,omitempty"`
}

func parsePath(p string) string {
	var filePath = p
	if filePath == "" {
		dirname, err := os.UserHomeDir()
		if err != nil {
			log.Fatalln("Failed to get user home dir:", err)
		}
		filePath = path.Join(dirname, ".chissl/profile.yaml")
	}
	return filePath
}

func fileExists(filePath string) bool {
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func NewClientConfig(profileConfigPath string) (*Config, error) {
	cfg := Config{Headers: http.Header{}}
	filePath := parsePath(profileConfigPath)

	if !fileExists(filePath) {
		if profileConfigPath != "" {
			return nil, errors.New("profile config file not found")
		}
		return &cfg, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("Failed to parse profile file at %s error: %v", filePath, err)
	}

	// Default values
	if cfg.MaxRetryCount == 0 {
		cfg.MaxRetryCount = -1
	}

	if cfg.KeepAlive == 0*time.Second {
		cfg.KeepAlive = 25 * time.Second
	}

	return &cfg, nil
}
