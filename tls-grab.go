package main

import "fmt"
import "os"
//import "net"
import "crypto/tls"

import "crypto/x509"
//import "crypto/rsa"
//import "crypto/dsa"
//import "crypto/ecdsa"

import "encoding/pem"

func main() {
  fmt.Println("google's public key:")

  conn, err := tls.Dial("tcp", "mail.google.com:443", &tls.Config{
    InsecureSkipVerify: true,
  })
  if err != nil {
    panic("failed to connect: " + err.Error())
  }
  conn.Close()
  tls_state := conn.ConnectionState()
  remote_certs := tls_state.PeerCertificates
  for cert_count, cert := range remote_certs {
    if (cert.BasicConstraintsValid) {
      if (cert.IsCA) {
        continue
      }
    } else {
      continue
    }

    fmt.Println("got cert #", cert_count)

    var pem_type string = ""

    switch cert.PublicKeyAlgorithm {
    case x509.RSA:
      fmt.Println("it's a RSA key")
      pem_type = "RSA PUBLIC KEY"

    case x509.ECDSA:
      fmt.Println("it's a ECDSA key")
      pem_type = "ECDSA PUBLIC KEY"

    case x509.DSA:
      fmt.Println("it's a DSA key")
      pem_type = "DSA PUBLIC KEY"

    default:
      fmt.Println("no clue")
    }

    pem_pubkey := &pem.Block{
      Type: pem_type,
      Bytes: cert.Raw,
    }
    pem.Encode(os.Stdout, pem_pubkey)
  }
  fmt.Println("didn't break?")
}
