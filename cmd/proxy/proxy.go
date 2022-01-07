package main

import (
   "bytes"
   "flag"
   "fmt"
   "github.com/89z/format/crypto"
   "io"
   "net"
   "net/http"
   "strconv"
)

type spyConn struct {
   SNI string
   filter bool
   net.Conn
}

func (s spyConn) Read(buf []byte) (int, error) {
   num, err := s.Conn.Read(buf)
   if !s.filter || bytes.Contains(buf, []byte(s.SNI)) {
      hello, err := crypto.ParseTLS(buf[:num])
      if err == nil {
         ja3, err := crypto.FormatJA3(hello)
         if err == nil {
            fmt.Printf("%q\n", buf[:num])
            fmt.Print("\t", ja3, "\n")
            fmt.Print("\t", crypto.Fingerprint(ja3), "\n")
         }
      }
   }
   return num, err
}

func (s spyConn) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
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
      s.Conn = clientConn
      io.Copy(targetConn, s)
   }
}

func main() {
   var (
      addr = ":8080"
      handler spyConn
   )
   flag.BoolVar(&handler.filter, "f", false, "filter")
   flag.StringVar(&handler.SNI, "s", "clientservices.googleapis.com", "SNI")
   flag.Parse()
   fmt.Println(addr)
   http.ListenAndServe(addr, handler)
}
