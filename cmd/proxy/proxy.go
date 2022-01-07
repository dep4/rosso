package main

import (
   "bytes"
   "fmt"
   "github.com/89z/format/crypto"
   "io"
   "net"
   "net/http"
   "strconv"
)

func (s spyConn) Read(buf []byte) (int, error) {
   num, err := s.Conn.Read(buf)
   if bytes.Contains(buf, []byte("clientservices.googleapis.com")) {
      hello, err := crypto.ParseTLS(buf[:num])
      if err == nil {
         ja3, err := crypto.FormatJA3(hello)
         if err == nil {
            fmt.Printf("%q\n\t%v\n\t%v\n", buf[:num], ja3, crypto.Fingerprint(ja3))
         }
      }
   }
   return num, err
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
