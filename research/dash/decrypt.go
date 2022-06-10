package decrypt

import (
   "github.com/edgeware/mp4ff/mp4"
   "io"
)

func decryptInit(r io.Reader, w io.Writer) error {
   inMp4, err := mp4.DecodeFile(r)
   if err != nil {
      return err
   }
   for _, trak := range inMp4.Init.Moov.Traks {
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
   return inMp4.Init.Encode(w)
}

func decryptSegment(r io.Reader, key []byte, w io.Writer) error {
   file, err := mp4.DecodeFile(r)
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
      err := seg.Encode(w)
      if err != nil {
         return err
      }
   }
   return nil
}
