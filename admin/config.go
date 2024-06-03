package chadmin

import (
	"errors"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

// Config represents a client configuration
type Config struct {
	//ProfilePath      string
	Fingerprint      string        `yaml:"fingerprint,omitempty"`
	Auth             string        `yaml:"auth,omitempty"`
	Username         string        `yaml:"username,omitempty"`
	Password         string        `yaml:"password,omitempty"`
	MaxRetryCount    int           `yaml:"max-retry-count,omitempty"`
	MaxRetryInterval time.Duration `yaml:"max-retry-interval,omitempty"`
	Server           string        `yaml:"server,omitempty"`
	Proxy            string        `yaml:"proxy,omitempty"`
	Headers          http.Header   `yaml:"headers,omitempty"`
	TLS              TLSConfig     `yaml:"tls,omitempty"`
	Verbose          bool          `yaml:"verbose,omitempty"`
}

// TLSConfig for a Client
type TLSConfig struct {
	SkipVerify bool   `yaml:"tls-skip-verify,omitempty"`
	CA         string `yaml:"tls-ca,omitempty"`
	Cert       string `yaml:"tls-cert,omitempty"`
	Key        string `yaml:"tls-key,omitempty"`
	ServerName string `yaml:"hostname,omitempty"`
}

func NewClientConfig(profileConfigPath string, cfg *Config) (*Config, error) {
	filePath := parsePath(profileConfigPath)

	if !fileExists(filePath) {
		if profileConfigPath != "" {
			return nil, errors.New("profile config file not found")
		}
		return cfg, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("Failed to parse profile file at %s error: %v", filePath, err)
	}

	cfg.Username, cfg.Password, err = parseCredentials(cfg.Auth)
	if err != nil {
		return nil, err
	}

	// Default values
	if cfg.MaxRetryCount == 0 {
		cfg.MaxRetryCount = -1
	}

	return cfg, nil
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

func parseCredentials(input string) (string, string, error) {
	// Check if the input is empty
	if input == "" {
		return "", "", errors.New("input cannot be empty")
	}

	// Split the input string by the colon
	parts := strings.Split(input, ":")

	// Check if there are exactly two parts
	if len(parts) != 2 {
		return "", "", errors.New("input must be in the format USERNAME:PASS")
	}

	// Extract username and password
	username := parts[0]
	password := parts[1]

	// Check if username or password is empty
	if username == "" || password == "" {
		return "", "", errors.New("username and password cannot be empty")
	}

	return username, password, nil
}

type RegexList struct {
	expressions []*regexp.Regexp
}

// String returns the string representation of the RegexList.
func (r *RegexList) String() string {
	exprs := make([]string, len(r.expressions))
	for i, expr := range r.expressions {
		exprs[i] = expr.String()
	}
	return strings.Join(exprs, ",")
}

// Set parses and sets the regex expressions from a comma-separated string.
func (r *RegexList) Set(value string) error {
	exprs := strings.Split(value, ",")
	r.expressions = make([]*regexp.Regexp, len(exprs))
	for i, expr := range exprs {
		re, err := regexp.Compile(expr)
		if err != nil {
			return err
		}
		r.expressions[i] = re
	}
	return nil
}
