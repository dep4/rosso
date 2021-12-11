package main

import (
   "crypto/tls"
   "fmt"
   "log"
)

type spyConn struct {
   *tls.Conn
}

func (s spyConn) Read(p []byte) (int, error) {
   n, err := s.Conn.Read(p)
   fmt.Printf("%q\n", p)
   return n, err
}

func main() {
    cer, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
    if err != nil {
        log.Println(err)
        return
    }
    config := &tls.Config{Certificates: []tls.Certificate{cer}}
    ln, err := tls.Listen("tcp", ":8080", config) 
    if err != nil {
        log.Println(err)
        return
    }
    defer ln.Close()
    for {
        conn, err := ln.Accept()
        if err != nil {
            log.Println(err)
            continue
        }
      defer conn.Close()
      var tConn spyConn
      tConn.Conn = conn.(*tls.Conn)
      fmt.Println("BEGIN")
      if err := tConn.Handshake(); err != nil {
         fmt.Println(err)
      }
      fmt.Println("END")
   }
}
