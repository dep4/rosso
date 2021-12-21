package main

import (
   "fmt"
   "github.com/89z/parse/crypto"
   "io"
   "net"
   "net/http"
   "time"
   //"crypto/tls"
)

func main() {
   server := http.Server{
      Addr: ":8080",
      Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
         if r.Method == http.MethodConnect {
            fmt.Println("handleTunneling")
            handleTunneling(w, r)
         }
      }),
   } 
   //err := server.ListenAndServeTLS("cert.pem", "key.pem")
   err := server.ListenAndServe()
   if err != nil {
      panic(err)
   }
}

func handleTunneling(w http.ResponseWriter, r *http.Request) {
   // net.Conn r
   dest_conn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
   if err != nil {
      fmt.Println(err)
      http.Error(w, err.Error(), http.StatusServiceUnavailable)
      return
   }
   w.WriteHeader(http.StatusOK)
   hijacker, ok := w.(http.Hijacker)
   if !ok {
      fmt.Println("Hijacker fail")
      http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
      return
   }
   // net.Conn w
   client_conn, _, err := hijacker.Hijack()
   if err != nil {
      fmt.Println(err)
      http.Error(w, err.Error(), http.StatusServiceUnavailable)
   }
   /*
   fail
   spy := spyConn{client_conn}
   go transfer(dest_conn, spy)
   go transfer(spy, dest_conn)
   */
   spy := spyConn{dest_conn}
   go transfer(spy, client_conn)
   go transfer(client_conn, spy)
}

// "\x16\x03\x01\x02\x00\x01\x00\x01\xfc"
func (s spyConn) Read(p []byte) (int, error) {
   n, err := s.Conn.Read(p)
   if p[0] == 0x16 {
      hand := crypto.Handshakes(p)[0]
      hello, err := crypto.ParseHandshake(hand)
      if err != nil {
         fmt.Println(err)
      } else {
         ja3, err := hello.FormatJA3()
         if err != nil {
            fmt.Println(err)
         } else {
            fmt.Println(ja3)
         }
      }
   }
   return n, err
}

type spyConn struct {
   net.Conn
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
    defer destination.Close()
    defer source.Close()
    io.Copy(destination, source)
}
