package main

import (
   "encoding/xml"
   "fmt"
   "github.com/89z/std/dash"
   "github.com/89z/std/mp4"
   "github.com/89z/std/os"
   "io"
)

func (f flags) DASH() error {
   res, err := amc.Client.Redirect(nil).Get(f.address)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   var media dash.Media
   if err := xml.NewDecoder(res.Body).Decode(&media); err != nil {
      return err
   }
   reps := media.Representations().Video()
   if f.bandwidth_video >= 1 {
      rep := reps.Get_Bandwidth(f.bandwidth_video)
      if f.info {
         for _, each := range reps {
            if each.Bandwidth == rep.Bandwidth {
               fmt.Print("!")
            }
            fmt.Println(each)
         }
      } else {
         var key []byte
         if source.Key_Systems != nil {
            key, err = f.key(play, rep.ContentProtection.Default_KID)
            if err != nil {
               return err
            }
         }
         return download(rep, key, play.Base())
      }
   }
   return nil
}


func download(rep *dash.Representation, key []byte, base string) error {
   file, err := os.Create(base + rep.Ext())
   if err != nil {
      return err
   }
   defer file.Close()
   res, err := amc.Client.Redirect(nil).Get(rep.Initialization())
   if err != nil {
      return err
   }
   defer res.Body.Close()
   media := rep.Media()
   pro := os.Progress_Chunks(file, len(media))
   dec := mp4.New_Decrypt(pro)
   if err := dec.Init(res.Body); err != nil {
      return err
   }
   for _, addr := range media {
      res, err := amc.Client.Redirect(nil).Level(0).Get(addr)
      if err != nil {
         return err
      }
      pro.Add_Chunk(res.ContentLength)
      if key != nil {
         err = dec.Segment(res.Body, key)
      } else {
         _, err = io.Copy(pro, res.Body)
      }
      if err != nil {
         return err
      }
      if err := res.Body.Close(); err != nil {
         return err
      }
   }
   return nil
}
