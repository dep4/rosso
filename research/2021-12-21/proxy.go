package main

import (
   "fmt"
   "github.com/89z/parse/crypto"
   "io"
   "net"
   "net/http"
   "strconv"
   "time"
)

const defaultTargetConnectTimeout   = 5 * time.Second

type proxy struct{}

func main() {
   serve := http.Server{
      Addr:       ":8080",
      ReadTimeout:  60 * time.Second,
      WriteTimeout: 60 * time.Second,
      Handler:    &proxy{},
   }
   fmt.Println("runServer")
   serve.ListenAndServe()
}

type spyConn struct {
   net.Conn
}

// "\x16\x03\x01\x02\x00\x01\x00\x01\xfc"
// 769,49195-49196-52393-49199-49200-52392-158-159-49161-49162-49171-49172-51-57-156-157-47-53,65281-0-23-35-13-16-11-10,23,0
func (s spyConn) Read(p []byte) (int, error) {
   n, err := s.Conn.Read(p)
   if p[0] == 0x16 {
      fmt.Println("Handshake")
      for _, hand := range crypto.Handshakes(p) {
         hello, err := crypto.ParseHandshake(hand)
         if err == nil {
            ja3, err := hello.FormatJA3()
            if err == nil {
               fmt.Println(ja3)
            }
         }
      }
   }
   return n, err
}

func (proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
   if req.URL.Host == "" {
      req.URL.Host = req.Host
   }
   if req.Method == http.MethodConnect {
      proxyTunnel(req, rw)
   }
}

func proxyTunnel(req *http.Request, rw http.ResponseWriter) {
   hijacker, ok := rw.(http.Hijacker)
   if !ok {
      fmt.Println("hijacker is not supported")
      return
   }
   clientConn, _, err := hijacker.Hijack()
   if err != nil {
      fmt.Printf("hijacker failed: %s", err)
      return
   }
   defer clientConn.Close()
   targetConn, err := net.DialTimeout("tcp", req.URL.Host, defaultTargetConnectTimeout)
   if err != nil {
      fmt.Printf("proxyTunnel %s dial remote server failed: %s", req.URL, err)
      return
   }
   defer targetConn.Close()
   buf := []byte("HTTP/1.1 ")
   buf = strconv.AppendInt(buf, http.StatusOK, 10)
   buf = append(buf, " Connection established\r\n\r\n"...)
   _, err = clientConn.Write(buf)
   if err != nil {
      fmt.Printf("proxyTunnel %s write message failed: %s", req.URL.Host, err)
      return
   }
   go func() {
      _, err1 := io.Copy(clientConn, targetConn)
      if err1 != nil {
         fmt.Println("io.Copy write clientConn failed:", err1)
      }
      clientConn.Close()
      targetConn.Close()
   }()
   spy := spyConn{clientConn}
   _, err2 := io.Copy(targetConn, spy)
   if err2 != nil {
      fmt.Println("io.Copy write targetConn failed:", err2)
   }
   targetConn.Close()
   clientConn.Close()
}
