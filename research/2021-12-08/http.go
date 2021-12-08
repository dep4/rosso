package main
 
import (
   "net/http"
   "net/http/httputil"
   "os"
)
 
func main() {
   http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
      buf, err := httputil.DumpRequest(req, true)
      if err != nil {
         panic(err)
      }
      os.Stdout.Write(buf)
   })
   http.ListenAndServe(":8080", nil)
}
