package dash

import (
   "github.com/edgeware/mp4ff/mp4"
   "io"
)

func decryptMP4withCenc(r io.Reader, key []byte, w io.Writer) error {
   inMp4, err := mp4.DecodeFile(r)
   if err != nil {
      return err
   }
   for _, seg := range inMp4.Segments {
      for _, frag := range seg.Fragments {
         moof := frag.Moof
         for _, traf := range moof.Trafs {
            samples, err := frag.GetFullSamples(nil)
            if err != nil {
               return err
            }
            for i := range samples {
               encSample := samples[i].Data
               var iv []byte
               if len(traf.Senc.IVs[i]) == 8 {
                  iv = make([]byte, 0, 16)
                  iv = append(iv, traf.Senc.IVs[i]...)
                  iv = append(iv, []byte{0, 0, 0, 0, 0, 0, 0, 0}...)
               } else {
                  iv = traf.Senc.IVs[i]
               }
               var sub []mp4.SubSamplePattern
               if len(traf.Senc.SubSamples) != 0 {
                  sub = traf.Senc.SubSamples[i]
               }
               dec, err := mp4.DecryptSampleCenc(encSample, key, iv, sub)
               if err != nil {
                  return err
               }
               copy(samples[i].Data, dec)
            }
            traf.RemoveEncryptionBoxes()
         }
         moof.RemovePsshs()
      }
      // fix jerk between fragments
      seg.Sidx = nil
      err := seg.Encode(w)
      if err != nil {
         return err
      }
   }
   return nil
}
