package dash

import (
   "github.com/edgeware/mp4ff/mp4"
   "io"
)

func Decrypt(dst io.Writer, src io.Reader, key []byte) error {
   file, err := mp4.DecodeFile(src)
   if err != nil {
      return err
   }
   for _, seg := range file.Segments {
      for _, frag := range seg.Fragments {
         for _, traf := range frag.Moof.Trafs {
            samples, err := frag.GetFullSamples(nil)
            if err != nil {
               return err
            }
            for i, sample := range samples {
               var iv []byte
               // this needs its own line so that the bytes are copied
               iv = append(iv, traf.Senc.IVs[i]...)
               iv = append(iv, 0, 0, 0, 0, 0, 0, 0, 0)
               var sub []mp4.SubSamplePattern
               if len(traf.Senc.SubSamples) > i {
                  sub = traf.Senc.SubSamples[i]
               }
               dec, err := mp4.DecryptSampleCenc(sample.Data, key, iv, sub)
               if err != nil {
                  return err
               }
               copy(sample.Data, dec)
            }
            // required for playback
            traf.RemoveEncryptionBoxes()
         }
         // fast start
         frag.Moof.RemovePsshs()
      }
      // fix jerk between fragments
      seg.Sidx = nil
      err := seg.Encode(dst)
      if err != nil {
         return err
      }
   }
   return nil
}

// Need for Mozilla Firefox and VLC media player
func DecryptInit(dst io.Writer, src io.Reader) error {
   file, err := mp4.DecodeFile(src)
   if err != nil {
      return err
   }
   for _, trak := range file.Init.Moov.Traks {
      for _, child := range trak.Mdia.Minf.Stbl.Stsd.Children {
         switch child.Type() {
         case "enca":
            _, err = child.(*mp4.AudioSampleEntryBox).RemoveEncryption()
         case "encv":
            _, err = child.(*mp4.VisualSampleEntryBox).RemoveEncryption()
         }
         if err != nil {
            return err
         }
      }
   }
   return file.Init.Encode(dst)
}
