package dash

import (
   "github.com/edgeware/mp4ff/mp4"
   "io"
)

func Decrypt(w io.Writer, r io.Reader, key []byte) error {
   file, err := mp4.DecodeFile(r)
   if err != nil {
      return err
   }
   if err := file.Init.Encode(w); err != nil {
      return err
   }
   for _, seg := range file.Segments {
      for _, frag := range seg.Fragments {
         for _, traf := range frag.Moof.Trafs {
            samples, err := frag.GetFullSamples(nil)
            if err != nil {
               return err
            }
            for i, sam := range samples {
               subSample := traf.Senc.SubSamples[i]
               iv := append(traf.Senc.IVs[i], 0, 0, 0, 0, 0, 0, 0, 0)
               dec, err := mp4.DecryptSampleCenc(sam.Data, key, iv, subSample)
               if err != nil {
                  return err
               }
               copy(sam.Data, dec)
            }
         }
      }
      err := seg.Encode(w)
      if err != nil {
         return err
      }
   }
   return nil
}