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
   if err := inMp4.Init.Encode(w); err != nil {
      return err
   }
   for _, seg := range inMp4.Segments {
      for _, frag := range seg.Fragments {
         for _, traf := range frag.Moof.Trafs {
            samples, err := frag.GetFullSamples(nil)
            if err != nil {
               return err
            }
            for i := range samples {
               encSample := samples[i].Data
               var iv []byte
               if len(traf.Senc.IVs[i]) == 8 {
                  iv = append(iv, traf.Senc.IVs[i]...)
                  iv = append(iv, []byte{0, 0, 0, 0, 0, 0, 0, 0}...)
               } else {
                  iv = traf.Senc.IVs[i]
               }
               subSamplePatterns := traf.Senc.SubSamples[i]
               decryptedSample, err := mp4.DecryptSampleCenc(encSample, key, iv, subSamplePatterns)
               if err != nil {
                  return err
               }
               copy(samples[i].Data, decryptedSample)
            }
            traf.RemoveEncryptionBoxes()
         }
      }
      if seg.Sidx != nil {
         seg.Sidx = nil // drop sidx inside segment, since not modified properly
      }
      err := seg.Encode(w)
      if err != nil {
         return err
      }
   }
   return nil
}

