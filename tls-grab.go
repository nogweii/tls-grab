package main

import (
	"crypto/sha256"
	"crypto/tls"
	"encoding/pem"
	"flag"
	"fmt"
	"os"
	"strings"
	"path"
	"log"
)

var (
	server      = flag.String("server", "", "The TLS host (name or IP) to connect to")
	port        = flag.Int("port", 443, "The port the TLS service is running on")
	network     = flag.String("net", "tcp", "Connect to this kind of network: tcp, udp, or unix")
	ipv4only    = flag.Bool("4", false, "Only connect via IPv4")
	ipv6only    = flag.Bool("6", false, "Only connect via IPv6")
	verify      = flag.Bool("verify", false, "Verify the provided certificate against trusted CAs")
	fingerprint = flag.Bool("fingerprint", false, "Print the SHA-256 fingerprint instead of the certificate")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			"Usage: %s [-4|-6] [-net=tcp|udp|unix] [-verify] [-port PORT] -server HOST\n",
			path.Base(os.Args[0]),
		)
		fmt.Fprintf(os.Stderr, "A simple utility to grab just a server's public certificate or fingerprint.\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}

	flag.Parse()

	if *server == "" {
		flag.Usage()
		log.Fatal("Need to specify a server.")
	}

	if *network != "tcp" && *network != "udp" && *network != "unix" {
		flag.Usage()
		log.Fatal("Unknown kind of network type! Try tcp.")
	}

	var networkSuffix string = ""
	if *ipv4only && *ipv6only {
		log.Print("Specifying both -4 & -6 is redundant.")
	} else if *ipv4only {
		networkSuffix = "4"
	} else if *ipv6only {
		networkSuffix = "6"
	}

	targetHost := fmt.Sprintf("%s:%d", *server, *port)
	networkType := fmt.Sprintf("%s%s", *network, networkSuffix)

	// Merely connect to the service, do the full handshake and then immediately
	// close the connection. All the information we care about is given in the
	// handshake.
	conn, err := tls.Dial(networkType, targetHost, &tls.Config{
		InsecureSkipVerify: !(*verify),
	})
	if err != nil {
		log.Fatalf("failed to connect: %s", err.Error())
	}
	conn.Close()

	tlsState := conn.ConnectionState()
	remoteCerts := tlsState.PeerCertificates
	for _, cert := range remoteCerts {
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
			// A certificate's fingerprint is simply the hash of the raw ASN.1 encoded
			// certificate. Super simple.
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
