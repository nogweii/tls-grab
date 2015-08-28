package main

import "fmt"
import "os"
import "crypto/tls"
import "crypto/x509"
import "encoding/pem"
import "flag"

func main() {

  var server = flag.String("server", "", "The TLS host (name or IP) to connect to")
  var port = flag.Int("port", 443, "The port the TLS service is running on")
  var network = flag.String("net", "tcp", "Connect to this kind of network: tcp, udp, or unix")
  var ipv4only = flag.Bool("4", false, "Only connect via IPv4")
  var ipv6only = flag.Bool("6", false, "Only connect via IPv6")
  flag.Parse()

  if (*server == "") {
    panic("Need to specify a server.")
  }

  if (*network != "tcp" && *network != "udp" && *network != "unix") {
    panic("Unknown kind of network type! Try tcp.")
  }

  var network_suffix string = ""
  if (*ipv4only && *ipv6only) {
    fmt.Println("Specifying both -4 & -6 is redundant")
  } else if (*ipv4only) {
    network_suffix = "4"
  } else if (*ipv6only) {
    network_suffix = "6"
  }

  target_host := fmt.Sprintf("%s:%d", *server, *port)
  network_type := fmt.Sprintf("%s%s", *network, network_suffix)

  conn, err := tls.Dial(network_type, target_host, &tls.Config{
    InsecureSkipVerify: true,
  })
  if err != nil {
    panic("failed to connect: " + err.Error())
  }
  conn.Close()
  tls_state := conn.ConnectionState()
  remote_certs := tls_state.PeerCertificates
  for _, cert := range remote_certs {
    if (cert.BasicConstraintsValid) {
      if (cert.IsCA) {
        continue
      }
    } else {
      continue
    }

    var pem_type string = ""

    switch cert.PublicKeyAlgorithm {
    case x509.RSA:
      pem_type = "RSA PUBLIC KEY"

    case x509.ECDSA:
      pem_type = "ECDSA PUBLIC KEY"

    case x509.DSA:
      pem_type = "DSA PUBLIC KEY"
    }

    pem_pubkey := &pem.Block{
      Type: pem_type,
      Bytes: cert.Raw,
    }
    pem.Encode(os.Stdout, pem_pubkey)
  }
}
