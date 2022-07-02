package main

import (
   "encoding/hex"
   "encoding/xml"
   "fmt"
   "github.com/89z/std/dash"
   "github.com/89z/std/mp4"
   "github.com/89z/std/os"
   "io"
)

func (f flags) DASH() error {
   res, err := client.Redirect(nil).Get(f.address)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   var media dash.Media
   if err := xml.NewDecoder(res.Body).Decode(&media); err != nil {
      return err
   }
   reps := media.Representations()
   var str stream
   str.base = "ignore"
   if f.key != "" {
      str.key, err = hex.DecodeString(f.key)
      if err != nil {
         return err
      }
   }
   str.Representations = reps.Audio()
   str.bandwidth = f.bandwidth_audio
   if err := f.download(str); err != nil {
      return err
   }
   str.Representations = reps.Video()
   str.bandwidth = f.bandwidth_video
   return f.download(str)
}

func (f flags) download(str stream) error {
   if str.bandwidth <= 0 {
      return nil
   }
   rep := str.Get_Bandwidth(str.bandwidth)
   if f.info {
      for _, each := range str.Representations {
         if each.Bandwidth == rep.Bandwidth {
            fmt.Print("!")
         }
         fmt.Println(each)
      }
      return nil
   }
   file, err := os.Create(str.base + rep.Ext())
   if err != nil {
      return err
   }
   defer file.Close()
   res, err := client.Redirect(nil).Get(rep.Initialization())
   if err != nil {
      return err
   }
   defer res.Body.Close()
   media := rep.Media()
   pro := os.Progress_Chunks(file, len(media))
   dec := mp4.New_Decrypt(pro)
   if str.key != nil {
      err = dec.Init(res.Body)
   } else {
      _, err = io.Copy(pro, res.Body)
   }
   if err != nil {
      return err
   }
   for _, addr := range media {
      res, err := client.Redirect(nil).Level(0).Get(addr)
      if err != nil {
         return err
      }
      pro.Add_Chunk(res.ContentLength)
      if str.key != nil {
         err = dec.Segment(res.Body, str.key)
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
