package chadmin

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/NextChapterSoftware/chissl/share/cio"
	"golang.org/x/sync/errgroup"
	"net/url"
	"os"
	"strings"
	"time"
)

// Admin represents an admin api client instance
type AdminClient struct {
	*cio.Logger
	config    *Config
	tlsConfig *tls.Config
	proxyURL  *url.URL
	server    string
	eg        *errgroup.Group
}

// NewAdminClient creates a new admin client instance
func NewAdminClient(c *Config) (*AdminClient, error) {
	//apply default scheme
	if !strings.HasPrefix(c.Server, "https://") {
		c.Server = "https://" + c.Server
	}
	if c.MaxRetryInterval < time.Second {
		c.MaxRetryInterval = 5 * time.Minute
	}
	u, err := url.Parse(c.Server)
	if err != nil {
		return nil, err
	}

	client := &AdminClient{
		Logger:    cio.NewLogger("client"),
		config:    c,
		server:    u.String(),
		tlsConfig: nil,
	}
	//set default log level
	client.Logger.Info = true
	//configure tls
	if u.Scheme == "wss" {
		tc := &tls.Config{}
		if c.TLS.ServerName != "" {
			tc.ServerName = c.TLS.ServerName
		}
		//certificate verification config
		if c.TLS.SkipVerify {
			client.Infof("TLS verification disabled")
			tc.InsecureSkipVerify = true
		} else if c.TLS.CA != "" {
			rootCAs := x509.NewCertPool()
			if b, err := os.ReadFile(c.TLS.CA); err != nil {
				return nil, fmt.Errorf("Failed to load file: %s", c.TLS.CA)
			} else if ok := rootCAs.AppendCertsFromPEM(b); !ok {
				return nil, fmt.Errorf("Failed to decode PEM: %s", c.TLS.CA)
			} else {
				client.Infof("TLS verification using CA %s", c.TLS.CA)
				tc.RootCAs = rootCAs
			}
		}
		//provide client cert and key pair for mtls
		if c.TLS.Cert != "" && c.TLS.Key != "" {
			c, err := tls.LoadX509KeyPair(c.TLS.Cert, c.TLS.Key)
			if err != nil {
				return nil, fmt.Errorf("Error loading client cert and key pair: %v", err)
			}
			tc.Certificates = []tls.Certificate{c}
		} else if c.TLS.Cert != "" || c.TLS.Key != "" {
			return nil, fmt.Errorf("Please specify client BOTH cert and key")
		}
		client.tlsConfig = tc
	}
	//validate remotes

	//outbound proxy
	if p := c.Proxy; p != "" {
		client.proxyURL, err = url.Parse(p)
		if err != nil {
			return nil, fmt.Errorf("Invalid proxy URL (%s)", err)
		}
	}

	return client, nil
}
