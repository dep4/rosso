package hls

import (
   "fmt"
   "io"
   "net/http"
   "os"
)

func read_response(res *http.Response) error {
   buf, err := io.ReadAll(res.Body)
   if err != nil {
      return err
   }
   return os.WriteFile("response.ts", buf, os.ModePerm)
}

func read_file(res *http.Response) error {
   file, err := os.Create("file")
   if err != nil {
      return err
   }
   if _, err := file.ReadFrom(res.Body); err != nil {
      return err
   }
   if err := file.Close(); err != nil {
      return err
   }
   file, err = os.Open("file")
   if err != nil {
      return err
   }
   defer file.Close()
   var b [9]byte
   if _, err := file.Read(b[:]); err != nil {
      return err
   }
   fmt.Println(b)
   return nil
}
