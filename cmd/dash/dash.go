package main

import (
   "encoding/xml"
   "flag"
   "fmt"
   "github.com/89z/mech/amc"
   "github.com/89z/mech/widevine"
   "github.com/89z/std/dash"
   "github.com/89z/std/mp4"
   "github.com/89z/std/os"
   "io"
   "path/filepath"
)

type flags struct {
   address string
   bandwidth_audio int
   bandwidth_video int
   info bool
   key string
}

func main() {
   var f flags
   flag.StringVar(&f.address, "a", "", "address")
   flag.IntVar(&f.bandwidth_video, "f", 1_999_999, "video bandwidth")
   flag.IntVar(&f.bandwidth_audio, "g", 127_000, "audio bandwidth")
   flag.BoolVar(&f.info, "i", false, "information")
   flag.StringVar(&f.key, "k", "", "key")
   flag.Parse()
   if f.address != "" {
      err := f.DASH()
      if err != nil {
         panic(err)
      }
   } else {
      flag.Usage()
   }
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

func (f *flags) key(p widevine.Poster, raw_key_id string) ([]byte, error) {
   private_key, err := os.ReadFile(f.private_key)
   if err != nil {
      return nil, err
   }
   client_ID, err := os.ReadFile(f.client_ID)
   if err != nil {
      return nil, err
   }
   key_ID, err := widevine.Key_ID(raw_key_id)
   if err != nil {
      return nil, err
   }
   mod, err := widevine.New_Module(private_key, client_ID, key_ID)
   if err != nil {
      return nil, err
   }
   keys, err := mod.Post(p)
   if err != nil {
      return nil, err
   }
   return keys.Content().Key, nil
}
func (f flags) DASH() error {
   home, err := os.UserHomeDir()
   if err != nil {
      return err
   }
   auth, err := amc.Open_Auth(home + "/mech/amc.json")
   if err != nil {
      return err
   }
   if err := auth.Refresh(); err != nil {
      return err
   }
   if err := auth.Create(home + "/mech/amc.json"); err != nil {
      return err
   }
   if f.nid == 0 {
      f.nid, err = amc.Get_NID(f.address)
      if err != nil {
         return err
      }
   }
   play, err := auth.Playback(f.nid)
   if err != nil {
      return err
   }
   source := play.Source()
   res, err := amc.Client.Redirect(nil).Get(source.Src)
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


