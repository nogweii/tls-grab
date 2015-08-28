package main

import "fmt"
//import "net"
import "crypto/tls"

func main() {
  fmt.Println("public key goes where?")
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
  for cert_count := range remote_certs {
    fmt.Println("got cert #", cert_count)
  }
  fmt.Println("didn't break?")
}
