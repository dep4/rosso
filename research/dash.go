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
   inMp4.Init.Encode(w)
   for _, seg := range inMp4.Segments {
      for _, frag := range seg.Fragments {
         for _, traf := range frag.Moof.Trafs {
            samples, err := frag.GetFullSamples(nil)
            if err != nil {
               return err
            }
            for i, sam := range samples {
               var iv []byte
               if len(traf.Senc.IVs[i]) == 8 {
                  iv = append(iv, traf.Senc.IVs[i]...)
                  iv = append(iv, 0, 0, 0, 0, 0, 0, 0, 0)
               } else {
                  iv = traf.Senc.IVs[i]
               }
               subSamplePatterns := traf.Senc.SubSamples[i]
               dec, err := mp4.DecryptSampleCenc(sam.Data, key, iv, subSamplePatterns)
               if err != nil {
                  return err
               }
               copy(samples[i].Data, dec)
            }
         }
      }
      seg.Encode(w)
   }
   return nil
}
