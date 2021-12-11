package main

import (
   "crypto/tls"
   "fmt"
   "net"
   "net/http"
)

func main() {
   http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
      w.WriteHeader(200)
   })
   // net.Listener
   ln, err := net.Listen("tcp", ":8080")
   if err != nil {
      panic(err)
   }
   defer ln.Close()
   for {
      // net.Conn
      conn, err := ln.Accept()
      if err != nil {
         fmt.Println(err)
         continue
      }
      defer conn.Close()
      rConn := spyConn{conn}
      var con tls.Config
      // tls.Conn
      tConn := tls.Server(rConn, &con)
      if err := tConn.Handshake(); err != nil {
         fmt.Println(err)
      }
   }
}

type spyConn struct {
   net.Conn
}

func (s spyConn) Read(p []byte) (int, error) {
   n, err := s.Conn.Read(p)
   fmt.Printf("%q\n", p)
   return n, err
}
