package main

import (
   "fmt"
   "github.com/89z/format/crypto"
   "io"
   "net"
   "net/http"
   "strconv"
)

func (s spyConn) Read(p []byte) (int, error) {
   n, err := s.Conn.Read(p)
   if hello, err := crypto.ParseTLS(p[:n]); err == nil {
      ja3, err := crypto.FormatJA3(hello)
      if err == nil {
         fmt.Print(crypto.Fingerprint(ja3), ":\n", ja3, "\n")
      }
   }
   return n, err
}

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

func main() {
   var (
      addr = ":8080"
      handler proxy
   )
   fmt.Println(addr)
   http.ListenAndServe(addr, handler)
}
