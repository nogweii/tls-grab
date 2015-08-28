package main

import "fmt"
//import "os"
//import "crypto/tls"
//import "crypto/x509"
//import "encoding/pem"
import "flag"

func main() {

  var server = flag.String("server", "mail.google.com", "The TLS host (name or IP) to connect to")
  var port = flag.Int("port", 443, "The port the TLS service is running on")
  flag.Parse()

  target_host := fmt.Sprintf("%s:%d", &server, &port)

  fmt.Println(target_host)

  /*
  conn, err := tls.Dial("tcp", target_host, &tls.Config{
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
  */
}
