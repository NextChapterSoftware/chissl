package settings

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

type Remote struct {
	UserAddress            string
	LocalHost, LocalPort   string
	RemoteHost, RemotePort string
	Reverse                bool
}

func validatePorts(port string) (int, error) {
	p, err := strconv.Atoi(port)
	if err != nil || p < 1 || p > 65535 {
		return 0, errors.New("invalid port number")
	}
	return p, nil
}

func validateHost(host string) error {
	if len(host) == 0 {
		return nil // Host is optional
	}

	// Check if the host is a valid IPv4 address
	if net.ParseIP(host) != nil {
		return nil
	}

	// Validate hostname according to RFC 1123
	if strings.ContainsAny(host, " !@#$%^&*()=+[]{}|;:'\",<>/?\\") {
		return errors.New("invalid hostname")
	}
	return nil
}

func DecodeRemote(s string) (*Remote, error) {
	parts := regexp.MustCompile(`^\s*(\d+)(?::([\w.-]+))?\s*->\s*(\d+)(?::([\w.-]+))?\s*$`).FindStringSubmatch(s)
	if len(parts) != 5 {
		return nil, errors.New("invalid remote format" + s)
	}

	// Local
	localPort := parts[1]
	localHost := parts[2]
	if localHost == "" {
		localHost = "0.0.0.0"
	}

	// Remote
	remotePort := parts[3]
	remoteHost := parts[4]
	if remoteHost == "" {
		remoteHost = "127.0.0.1"
	}

	// Validate ports
	if _, err := validatePorts(remotePort); err != nil {
		return nil, fmt.Errorf("invalid remote port: %v", err)
	}
	if _, err := validatePorts(localPort); err != nil {
		return nil, fmt.Errorf("invalid local port: %v", err)
	}

	// Validate remote host
	if err := validateHost(remoteHost); err != nil {
		return nil, fmt.Errorf("invalid remote host: %v", err)
	}

	// Validate local host
	if err := validateHost(localHost); err != nil {
		return nil, fmt.Errorf("invalid local host: %v", err)
	}

	r := &Remote{
		UserAddress: strings.TrimSpace(s),
		LocalHost:   localHost,
		LocalPort:   localPort,
		RemoteHost:  remoteHost,
		RemotePort:  remotePort,
		Reverse:     true,
	}
	return r, nil
}

var l4Proto = regexp.MustCompile(`(?i)\/(tcp|udp)$`)

// L4Proto extacts the layer-4 protocol from the given string
func L4Proto(s string) (head, proto string) {
	if l4Proto.MatchString(s) {
		l := len(s)
		return strings.ToLower(s[:l-4]), s[l-3:]
	}
	return s, ""
}

// implement Stringer
func (r Remote) String() string {
	sb := strings.Builder{}
	sb.WriteString(strings.TrimPrefix(r.Local(), "0.0.0.0:"))
	sb.WriteString("->")
	sb.WriteString(strings.TrimPrefix(r.Remote(), "127.0.0.1:"))
	return sb.String()
}

// Encode remote to a string
func (r Remote) Encode() string {
	return r.Local() + "->" + r.Remote()
}

// Local is the decodable local portion
func (r Remote) Local() string { return r.LocalHost + ":" + r.LocalPort }

// Remote is the decodable remote portion
func (r Remote) Remote() string {
	return r.RemoteHost + ":" + r.RemotePort
}

// UserAddr is checked when checking if a
// user has access to a given remote
func (r Remote) UserAddr() string {
	return r.UserAddress
}

// CanListen checks if the port can be listened on
func (r Remote) CanListen() bool {
	conn, err := net.Listen("tcp", r.Local())
	if err == nil {
		conn.Close()
		return true
	}
	return false
}

type Remotes []*Remote

// TODO: Legacy to be removed
func (rs Remotes) Reversed(reverse bool) Remotes {
	return rs
}

// Encode back into strings
func (rs Remotes) Encode() []string {
	s := make([]string, len(rs))
	for i, r := range rs {
		s[i] = r.Encode()
	}
	return s
}
