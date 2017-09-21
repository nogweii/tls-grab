# tls-grab

A simple utility to grab just a server's public certificate or fingerprint.

Verification is disabled by default. If you want that, pass `-verify`.

Licensed under the MIT license. See LICENSE for details.

## Install

Installation is simple, assuming you already have [go set up][goinstall]:

    go get -u github.com/evaryont/tls-grab

## Usage

```
Usage: tls-grab [-afVy] [-4|-6] [--net=tcp|udp|unix] [--port 443] [--host host] server[:port]
A simple utility to grab just a server's public certificate or fingerprint.

  -a, --all           Show every certificate, don't skip CA and intermediates
  -f, --fingerprint   Print the SHA-256 fingerprint instead of the certificate
  -s, --host host     Override what hostname to use when connecting to a server (for SNI)
  -4, --ipv4          Only connect via IPv4
  -6, --ipv6          Only connect via IPv6
  -n, --net string    Connect to this kind of network: tcp, udp, or unix (default "tcp")
  -p, --port int      The port the TLS service is running on (default 443)
  -V, --verbose       Print extra log messages
  -y, --verify        Verify the provided certificate against trusted CAs
      --version       Display the version
```

[goinstall]: https://golang.org/doc/install#install
