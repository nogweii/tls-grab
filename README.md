# tls-grab

A very simple tool to grab the public certificate or key of the remote server.

Verification is disabled by default. If you want that, pass `-verify`.

Licensed under the MIT license.

## Usage

```
Usage of tls-grab:
  -4    Only connect via IPv4
  -6    Only connect via IPv6
  -fingerprint
        Print the SHA-256 fingerprint instead of the certificate
  -net string
        Connect to this kind of network: tcp, udp, or unix (default "tcp")
  -port int
        The port the TLS service is running on (default 443)
  -server string
        The TLS host (name or IP) to connect to
  -verify
        Verify the provided certificate against trusted CAs
```
