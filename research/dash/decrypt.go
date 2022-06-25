package dash

import (
   "github.com/edgeware/mp4ff/mp4"
   "io"
)

type sinf_box struct {
   *mp4.SinfBox
}

func decrypt(init io.Reader, w io.Writer) (*sinf_box, error) {
   file, err := mp4.DecodeFile(init)
   if err != nil {
      return nil, err
   }
   var sinf sinf_box
   for _, trak := range file.Init.Moov.Traks {
      for _, child := range trak.Mdia.Minf.Stbl.Stsd.Children {
         switch child.Type() {
         case "encv":
            sinf.SinfBox, err = child.(*mp4.VisualSampleEntryBox).RemoveEncryption()
         case "enca":
            sinf.SinfBox, err = child.(*mp4.AudioSampleEntryBox).RemoveEncryption()
         }
         if err != nil {
            return nil, err
         }
      }
   }
   file.Init.Moov.RemovePsshs()
   if err := file.Init.Encode(w); err != nil {
      return nil, err
   }
   return &sinf, nil
}

func (s sinf_box) decrypt(segment io.Reader, key []byte, w io.Writer) error {
   file, err := mp4.DecodeFile(segment)
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
            tenc := s.SinfBox.Schi.Tenc
            for i, sample := range samples {
               var iv []byte
               if len(traf.Senc.IVs) == len(samples) {
                  if len(traf.Senc.IVs[i]) == 8 {
                     iv = append(iv, traf.Senc.IVs[i]...)
                     iv = append(iv, []byte{0, 0, 0, 0, 0, 0, 0, 0}...)
                  } else if len(traf.Senc.IVs) == len(samples) {
                     iv = traf.Senc.IVs[i]
                  }
               } else if tenc.DefaultConstantIV != nil {
                  iv = tenc.DefaultConstantIV
               }
               var subSamplePatterns []mp4.SubSamplePattern
               if len(traf.Senc.SubSamples) >= 1 {
                  subSamplePatterns = traf.Senc.SubSamples[i]
               }
               switch s.SinfBox.Schm.SchemeType {
               case "cenc":
                  err = mp4.DecryptSampleCenc(sample.Data, key, iv, subSamplePatterns)
               case "cbcs":
                  err = mp4.DecryptSampleCbcs(sample.Data, key, iv, subSamplePatterns, tenc)
               }
               if err != nil {
                  return err
               }
            }
            traf.RemoveEncryptionBoxes()
         }
         frag.Moof.RemovePsshs()
      }
      seg.Sidx = nil // drop sidx inside segment, since not modified properly
      err := seg.Encode(w)
      if err != nil {
         return err
      }
   }
   return nil
}
