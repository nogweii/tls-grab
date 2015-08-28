package main

import "fmt"
//import "net"
import "crypto/tls"

import "crypto/x509"
//import "crypto/rsa"
//import "crypto/dsa"
//import "crypto/ecdsa"

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
    fmt.Println("got cert #", cert_count)

    switch cert.PublicKeyAlgorithm {
    case x509.RSA:
      fmt.Println("it's a RSA key")

    case x509.ECDSA:
      fmt.Println("it's a ECDSA key")

    case x509.DSA:
      fmt.Println("it's a DSA key")

    default:
      fmt.Println("no clue")
    }
  }
  fmt.Println("didn't break?")
}
