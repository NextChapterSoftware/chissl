# chiSSL

## TL;DR

Jump to [examples](#example-usage)

## HTTPS Reverse Tunnel Made Easy!

chiSSL is essentially a light version of [Chisel](https://github.com/jpillora/chisel) exclusively focusing on creating **on-demand** HTTPS reverse tunnels.

### Features:
- Fully managed SSL certificates via Let's Encrypt
- [Payload inspection](#payload-inspection) for connections 
- Multiple concurrent HTTPS reverse tunnels 
- Client configuration using [yaml file](#client-usage)
- User management CLI (beta)
- User management REST API (beta)

![chissl Reverse HTTPS Tunnel](./assets/chissl.png)

## Requirements:

- Linux server with a public IP address
- Fully qualified domain name!

## Installation

### Homebrew

```bash
brew tap nextchaptersoftware/chissl https://github.com/NextChapterSoftware/chissl
brew install chissl
```


### Linux Installer (Server only)

```bash
# <domain_name> must be set to your server's fqdn 
# [port] is optional and will default to 443
# bash <(curl -s https://raw.githubusercontent.com/NextChapterSoftware/chissl/main/server_installer.sh) <domain_name> [port]
# e.g 
bash <(curl -s https://raw.githubusercontent.com/NextChapterSoftware/chissl/main/server_installer.sh) your.domain.com
```

### Manual
Binaries for all supported platforms can be downloaded from [Releases](https://github.com/NextChapterSoftware/chissl/releases) page.

```bash
# sudo curl -L -o /usr/local/bin/chissl https://github.com/NextChapterSoftware/chissl/releases/download/v{VERSION}/chissl_{VERSION}_{OS}_{ARCH}
# e.g 
sudo curl -L -o /usr/local/bin/chissl https://github.com/NextChapterSoftware/chissl/releases/download/v1.2/chissl_1.2_linux_amd64
sudo chmod +x /usr/local/bin/chissl
```

<h2 id="example-usage">
Examples 
</h2>

### Server

```bash 
# Fully automated Let's Encrypt certificate
chissl server --tls-domain tunnel.example.com --auth Username:Password

# Bring your own certificate
chissl server --tls-key /path/to/your/privkey.pem \
                  --tls-cert /path/to/your/fullchain.pem \
                  --auth Username:Password 
```

### Client

```bash 
# For a single port 
chissl client --auth Username:Password  https://tunnel.example.com "PortOnServer->PortOnYourMachine"

# It can work with multiple port mappings too! 
chissl client --auth Username:Password  https://tunnel.example.com \
      "Port1_OnServer->PortX_OnYourMachine" "Port2_OnServer->PortY_OnYourMachine" "Port3_OnServer->PortZ_OnYourMachine"
```

## Authentication

Using the --authfile option, the server may optionally provide a user.json configuration file to create a list of accepted users. The client then authenticates using the --auth option. See [users.json](example/users.json) for an example authentication configuration file. See the --help above for more information.

Internally, this is done using the Password authentication method provided by SSH. Learn more about crypto/ssh here http://blog.gopheracademy.com/go-and-ssh/.

<h2 id="payload-inspection">
Payload Inspection 
</h2>
Setting client logging to debug (`-v`) will output data sent/received along with stats for each connection. 

Example output:
```plain
2024/05/20 23:46:54 client: tun: conn#1: Open [1/1]
2024/05/20 23:46:54 client: tun: conn#1:
================== Host: 127.0.0.1:8000 : Read ==================
GET / HTTP/1.1
Host: tunnel.example.com:8098
User-Agent: curl/8.4.0
Accept: */*

2024/05/20 23:46:54 client: tun: conn#1:
================== Host: 127.0.0.1:8000 : Write ==================
HTTP/1.0 200 OK
Server: SimpleHTTP/0.6 Python/3.11.3
Date: Tue, 21 May 2024 06:46:54 GMT
Content-type: text/html; charset=utf-8
Content-Length: 187

2024/05/20 23:46:54 client: tun: conn#1:
================== Host: 127.0.0.1:8000 : Write ==================
<!DOCTYPE HTML>
<html lang="en">
<head>
<meta charset="utf-8">
<title>Directory listing for /</title>
</head>
<body>
<h1>Directory listing for /</h1>
<hr>
<ul>
</ul>
<hr>
</body>
</html>
2024/05/20 23:46:54 client: tun: conn#1: sent 98B received 342B
2024/05/20 23:46:54 client: tun: conn#1: Close [0/1]
```

## Usage

```plain 
$ chissl --help

  Usage: chissl [command] [--help]

  Version: X.Y.Z

  Commands:
    server - runs chissl in server mode
    client - runs chissl in client mode

  Read more:
    https://github.com/NextChapterSoftware/chissl

```

### Server

```plain
Usage: chissl server [options]

  Options:

    --host, Defines the HTTP listening host – the network interface
    (defaults the environment variable HOST and falls back to 0.0.0.0).

    --port, -p, Defines the HTTP listening port (defaults to the environment
    variable PORT and fallsback to port 8080).

    --key, (deprecated use --keygen and --keyfile instead)
    An optional string to seed the generation of a ECDSA public
    and private key pair. All communications will be secured using this
    key pair. Share the subsequent fingerprint with clients to enable detection
    of man-in-the-middle attacks (defaults to the CHISEL_KEY environment
    variable, otherwise a new key is generate each run).

    --keygen, A path to write a newly generated PEM-encoded SSH private key file.
    If users depend on your --key fingerprint, you may also include your --key to
    output your existing key. Use - (dash) to output the generated key to stdout.

    --keyfile, An optional path to a PEM-encoded SSH private key. When
    this flag is set, the --key option is ignored, and the provided private key
    is used to secure all communications. (defaults to the CHISEL_KEY_FILE
    environment variable). Since ECDSA keys are short, you may also set keyfile
    to an inline base64 private key (e.g. chissl server --keygen - | base64).

    --authfile, An optional path to a users.json file. This file should
    be an object with users defined like:
      {
        "<user:pass>": ["<addr-regex>","<addr-regex>"]
      }
    when <user> connects, their <pass> will be verified and then
    each of the remote addresses will be compared against the list
    of address regular expressions for a match. Addresses will always 
    come in the form:
        "local-port:local-host->remote-port:remote-host" 
    This file will be automatically reloaded on change.

    --auth, An optional string representing a single user with full
    access, in the form of <user:pass>. It is equivalent to creating an
    authfile with {"<user:pass>": [""]}. If unset, it will use the
    environment variable AUTH.

    --keepalive, An optional keepalive interval. Since the underlying
    transport is HTTP, in many instances we'll be traversing through
    proxies, often these proxies will close idle connections. You must
    specify a time with a unit, for example '5s' or '2m'. Defaults
    to '25s' (set to 0s to disable).

    --tls-key, Enables TLS and provides optional path to a PEM-encoded
    TLS private key. When this flag is set, you must also set --tls-cert,
    and you cannot set --tls-domain.

    --tls-cert, Enables TLS and provides optional path to a PEM-encoded
    TLS certificate. When this flag is set, you must also set --tls-key,
    and you cannot set --tls-domain.

    --tls-domain, Enables TLS and automatically acquires a TLS key and
    certificate using LetsEncrypt. Setting --tls-domain requires port 443.
    You may specify multiple --tls-domain flags to serve multiple domains.
    The resulting files are cached in the "$HOME/.cache/chisel" directory.
    You can modify this path by setting the CHISEL_LE_CACHE variable,
    or disable caching by setting this variable to "-". You can optionally
    provide a certificate notification email by setting CHISEL_LE_EMAIL.

    --tls-ca, a path to a PEM encoded CA certificate bundle or a directory
    holding multiple PEM encode CA certificate bundle files, which is used to 
    validate client connections. The provided CA certificates will be used 
    instead of the system roots. This is commonly used to implement mutual-TLS.
```

<h3 id="client-usage">
Client 
</h3>

Clients can also be configured using a yaml file placed at `.chissl/profile.yaml` under current user's home directory or at a path supplied using `--profile` option. 

See [profile-example.yaml](example/profile-example.yaml) and [profile-example-full.yaml](example/profile-example-full.yaml).

```plain
Usage: chissl client [options] <server> <remote> [remote] [remote] ...

  <server> is the URL to the chissl server.

  <remote>s are remote connections tunneled through the server, each of
  which come in the form:
   	local-port:local-host->remote-port:remote-host
	
    ■ local-port (port on server) is required*.
    ■ local-host (interface on server) defaults to 0.0.0.0 (all interfaces).
    ■ remote-port is required*.
    ■ remote-host defaults to 127.0.0.1
    
    example remotes
      8080->80
      8080:0.0.0.0->80
      8080:0.0.0.0->80:127.0.0.1
      8089->80:neverssl.com

  Options:
  
    --profile, path to profile configuration yaml file. Defaults to 
    $HOME/chissl/profile.yaml. Profile yaml file allows users to 
    set all client arguments and configurations using a static file. 
	
	YAML Options:
	---
	fingerprint: "sample_fingerprint"
	auth: "user:password"
	keepalive: 30s
	max-retry-count: 10
	max-retry-interval: 2m
	server: "example.com"
	proxy: "http://admin:password@my-server.com:8081"
	remotes:
	  - 8089->80:neverssl.com
      - 8080->80
	headers:
	  Foo: ["Bar"]
	tls:
	  tls-skip-verify: true
	  tls-ca: "/path/to/ca"
	  tls-cert: "/path/to/cert"
	  tls-key: "/path/to/key"
	  hostname: "example.com"
	verbose: true

    --fingerprint, A *strongly recommended* fingerprint string
    to perform host-key validation against the server's public key.
    Fingerprint mismatches will close the connection.
    Fingerprints are generated by hashing the ECDSA public key using
    SHA256 and encoding the result in base64.
    Fingerprints must be 44 characters containing a trailing equals (=).

    --auth, An optional username and password (client authentication)
    in the form: "<user>:<pass>". These credentials are compared to
    the credentials inside the server's --authfile. defaults to the
    AUTH environment variable.

    --keepalive, An optional keepalive interval. Since the underlying
    transport is HTTP, in many instances we'll be traversing through
    proxies, often these proxies will close idle connections. You must
    specify a time with a unit, for example '5s' or '2m'. Defaults
    to '25s' (set to 0s to disable).

    --max-retry-count, Maximum number of times to retry before exiting.
    Defaults to unlimited.

    --max-retry-interval, Maximum wait time before retrying after a
    disconnection. Defaults to 5 minutes.

    --proxy, An optional HTTP CONNECT or SOCKS5 proxy which will be
    used to reach the chissl server. Authentication can be specified
    inside the URL.
    For example, http://admin:password@my-server.com:8081
            or: socks://admin:password@my-server.com:1080

    --hostname, Optionally set the 'Host' header (defaults to the host
    found in the server url).

    --sni, Override the ServerName when using TLS (defaults to the 
    hostname).

    --tls-ca, An optional root certificate bundle used to verify the
    chissl server. Only valid when connecting to the server with
    "https" or "wss". By default, the operating system CAs will be used.

    --tls-skip-verify, Skip server TLS certificate verification of
    chain and host name (if TLS is used for transport connections to
    server). If set, client accepts any TLS certificate presented by
    the server and any host name in that certificate. This only affects
    transport https (wss) connection. Chisel server's public key
    may be still verified (see --fingerprint) after inner connection
    is established.

    --tls-key, a path to a PEM encoded private key used for client 
    authentication (mutual-TLS).

    --tls-cert, a path to a PEM encoded certificate matching the provided 
    private key. The certificate must have client authentication 
    enabled (mutual-TLS)
```

### Admin CLI (beta)
Admin CLI docs can be found [here](./docs/admin_cli.md.md)

### User management REST API (beta)
User management REST API docs can be found [here](./docs/rest_api.md) 

## Roadmap 

- Used hashed/obfuscated passwords in server config file  
- Modernise CLI implementation 
- Improve verbose logging method 
- Fix the ability to set/forward default headers 