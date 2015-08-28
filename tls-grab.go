package main

import (
	"crypto/sha256"
	"crypto/tls"
	"encoding/pem"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {

	var server = flag.String("server", "", "The TLS host (name or IP) to connect to")
	var port = flag.Int("port", 443, "The port the TLS service is running on")
	var network = flag.String("net", "tcp", "Connect to this kind of network: tcp, udp, or unix")
	var ipv4only = flag.Bool("4", false, "Only connect via IPv4")
	var ipv6only = flag.Bool("6", false, "Only connect via IPv6")
	var verify = flag.Bool("verify", false, "Verify the provided certificate against trusted CAs")
	var fingerprint = flag.Bool("fingerprint", false, "Print the SHA-256 fingerprint instead of the certificate")
	flag.Parse()

	if *server == "" {
		fmt.Fprintln(os.Stderr, "Need to specify a server.")
		flag.Usage()
		os.Exit(2)
	}

	if *network != "tcp" && *network != "udp" && *network != "unix" {
		fmt.Fprintln(os.Stderr, "Unknown kind of network type! Try tcp.")
		flag.Usage()
		os.Exit(2)
	}

	var network_suffix string = ""
	if *ipv4only && *ipv6only {
		fmt.Println("Specifying both -4 & -6 is redundant")
	} else if *ipv4only {
		network_suffix = "4"
	} else if *ipv6only {
		network_suffix = "6"
	}

	target_host := fmt.Sprintf("%s:%d", *server, *port)
	network_type := fmt.Sprintf("%s%s", *network, network_suffix)

	// Merely connect to the service, do the full handshake and then immediately
	// close the connection. All the information we care about is given in the
	// handshake.
	conn, err := tls.Dial(network_type, target_host, &tls.Config{
		InsecureSkipVerify: !(*verify),
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to connect: "+err.Error())
		os.Exit(1)
	}
	conn.Close()

	tls_state := conn.ConnectionState()
	remote_certs := tls_state.PeerCertificates
	for _, cert := range remote_certs {
		// skip all of the intermediary & root CA certificates provided by the
		// server
		if cert.BasicConstraintsValid {
			if cert.IsCA {
				continue
			}
		} else {
			continue
		}

		if *fingerprint {
			fmt.Println(strings.Replace(fmt.Sprintf("% X", sha256.Sum256(cert.Raw)),
				" ", ":", -1))
		} else {
			pem.Encode(os.Stdout, &pem.Block{
				Type:  "CERTIFICATE",
				Bytes: cert.Raw,
			})
		}
	}
}
