package main

import (
   "fmt"
   "github.com/89z/parse/crypto"
   "io"
   "net"
   "net/http"
   "strconv"
)

type proxy struct{}

func (proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
   if req.Method == http.MethodConnect {
      hijacker, ok := rw.(http.Hijacker)
      if !ok {
         fmt.Println("interface is not http.Hijacker")
         return
      }
      clientConn, _, err := hijacker.Hijack()
      if err != nil {
         fmt.Println(err)
         return
      }
      defer clientConn.Close()
      targetConn, err := net.Dial("tcp", req.URL.Host)
      if err != nil {
         fmt.Println(err)
         return
      }
      defer targetConn.Close()
      buf := []byte("HTTP/1.1 ")
      buf = strconv.AppendInt(buf, http.StatusOK, 10)
      buf = append(buf, "\n\n"...)
      clientConn.Write(buf)
      spy := spyConn{clientConn}
      io.Copy(targetConn, spy)
   }
}

type spyConn struct {
   net.Conn
}

func (s spyConn) Read(p []byte) (int, error) {
   n, err := s.Conn.Read(p)
   if hello, err := crypto.ParseHandshake(p[:n]); err == nil {
      ja3, err := hello.FormatJA3()
      if err == nil {
         fmt.Println(ja3)
      }
   }
   return n, err
}

func main() {
   var hand proxy
   fmt.Println("ListenAndServe")
   http.ListenAndServe(":8080", hand)
}
