package main

import (
   "crypto/tls"
   "fmt"
   "net"
   "net/http"
)

func main() {
   tra := http.Transport{
      DialTLS: func(network, addr string) (net.Conn, error) {
         conn, err := net.Dial(network, addr)
         if err != nil {
            return nil, err
         }
         host, _, err := net.SplitHostPort(addr)
         if err != nil {
            return nil, err
         }
         config := &tls.Config{ServerName: host}
         tConn := tls.Client(conn, config)
         if err := tConn.Handshake(); err != nil {
            return nil, err
         }
         return tConn, nil
      },
   }
   req, err := http.NewRequest("HEAD", "http://example.com", nil)
   if err != nil {
      panic(err)
   }
   res, err := tra.RoundTrip(req)
   if err != nil {
      panic(err)
   }
   fmt.Printf("%+v\n", res)
}
