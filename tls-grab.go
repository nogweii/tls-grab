package main

import (
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
	"path"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

var (
	hostname    = flag.StringP("host", "s", "", "Override what `host`name to use when connecting to a server (for SNI)")
	port        = flag.IntP("port", "p", 443, "The port the TLS service is running on")
	network     = flag.StringP("net", "n", "tcp", "Connect to this kind of network: tcp, udp, or unix")
	ipv4only    = flag.BoolP("ipv4", "4", false, "Only connect via IPv4")
	ipv6only    = flag.BoolP("ipv6", "6", false, "Only connect via IPv6")
	verify      = flag.BoolP("verify", "y", false, "Verify the provided certificate against trusted CAs")
	fingerprint = flag.BoolP("fingerprint", "f", false, "Print the SHA-256 fingerprint instead of the certificate")
	verbose     = flag.BoolP("verbose", "V", false, "Print extra log messages")
	all_certs   = flag.BoolP("all", "a", false, "Show every certificate, don't skip CA and intermediates")
	version     = flag.Bool("version", false, "Display the version")
)

var Version = "2.0.0"

func displayCert(cert x509.Certificate) {
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

func main() {
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			"Usage: %s [-afVy] [-4|-6] [--net=tcp|udp|unix] [--port 443] [--host host] server[:port]\n",
			path.Base(os.Args[0]),
		)
		fmt.Fprintf(os.Stderr, "A simple utility to grab just a server's public certificate or fingerprint.\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}

	flag.Parse()

  log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
	})
	if *verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}

	if len(flag.Args()) == 0 {
		flag.Usage()
		log.Fatal("Missing server to connect to.")
	}

	if *network != "tcp" && *network != "udp" && *network != "unix" {
		flag.Usage()
		log.Fatal("Unknown kind of network type! Try tcp.")
	}

	var networkSuffix string = ""
	if *ipv4only && *ipv6only {
		log.Error("Specifying both -4 & -6 is redundant.")
	} else if *ipv4only {
		networkSuffix = "4"
	} else if *ipv6only {
		networkSuffix = "6"
	}

	server := flag.Args()[0]
	targetHost := fmt.Sprintf("%s:%d", server, *port)
	networkType := fmt.Sprintf("%s%s", *network, networkSuffix)

  log.WithFields(log.Fields{
    "host": targetHost,
    "network": networkType,
		"verify": *verify,
  }).Info("Connecting to host")

	// Merely connect to the service, do the full handshake and then immediately
	// close the connection. All the information we care about is given in the
	// handshake.
	conn, err := tls.Dial(networkType, targetHost, &tls.Config{
		InsecureSkipVerify: !(*verify),
		ServerName: *hostname,
	})
	if err != nil {
		log.Fatalf("failed to connect: %s", err.Error())
	}
	conn.Close()

	tlsState := conn.ConnectionState()
	remoteCerts := tlsState.PeerCertificates
	log.Infof("got %d certificates", len(remoteCerts))

	if *all_certs {
		// loop through and display every certificate found
		for _, cert := range remoteCerts {
			displayCert(*cert)
		}
	} else {
		// get the last certificate shown and present it
		//cert := remoteCerts[len(remoteCerts)-1]
		cert := remoteCerts[0]
		displayCert(*cert)
	}
}
