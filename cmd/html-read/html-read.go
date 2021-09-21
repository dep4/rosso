package main

import (
   "encoding/json"
   "fmt"
   "github.com/89z/parse/html"
   "net/http"
   "os"
)

func main() {
   if len(os.Args) != 2 {
      fmt.Println("media [URL]")
      return
   }
   addr := os.Args[1]
   fmt.Println("GET", addr)
   res, err := http.Get(addr)
   if err != nil {
      panic(err)
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      panic(res.Status)
   }
   lex := html.NewLexer(res.Body)
   // This is going to kill audio and video if the page is missing og:image.
   // However that is unlikely, so we will cross that bridge when we come to it.
   lex.NextAttr("property", "og:image")
   fmt.Println(lex.GetAttr("content"))
   // audio video
   for lex.NextAttr("type", "application/ld+json") {
      data := lex.Bytes()
      // audio
      var audio struct {
         ContentURL string
      }
      json.Unmarshal(data, &audio)
      if audio.ContentURL != "" {
         fmt.Println(audio.ContentURL)
      }
      // video
      var article struct {
         Video struct {
            ContentURL string
         }
      }
      json.Unmarshal(data, &article)
      if article.Video.ContentURL != "" {
         fmt.Println(article.Video.ContentURL)
      }
   }
}
