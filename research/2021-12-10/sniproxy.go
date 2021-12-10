package main

import (
   "bytes"
   "crypto/tls"
   "fmt"
   "io"
   "net"
   "time"
)

func main() {
   l, err := net.Listen("tcp", ":8080")
   if err != nil {
      panic(err)
   }
   for {
      conn, err := l.Accept()
      if err != nil {
         fmt.Println(err)
         continue
      }
      defer conn.Close()
      peekedBytes := new(bytes.Buffer)
      con := new(tls.Config)
      var rConn readOnlyConn
      rConn.reader = io.TeeReader(conn, peekedBytes)
      if err := tls.Server(rConn, con).Handshake(); err != nil {
         fmt.Println(err)
      }
   }
}

type readOnlyConn struct {
	reader io.Reader
}

func (conn readOnlyConn) Read(p []byte) (int, error)         {
   n, err := conn.reader.Read(p)
   fmt.Printf("%q\n", p)
   return n, err
}

func (conn readOnlyConn) Write(p []byte) (int, error)        { return 0, io.ErrClosedPipe }
func (conn readOnlyConn) Close() error                       { return nil }
func (conn readOnlyConn) LocalAddr() net.Addr                { return nil }
func (conn readOnlyConn) RemoteAddr() net.Addr               { return nil }
func (conn readOnlyConn) SetDeadline(t time.Time) error      { return nil }
func (conn readOnlyConn) SetReadDeadline(t time.Time) error  { return nil }
func (conn readOnlyConn) SetWriteDeadline(t time.Time) error { return nil }
