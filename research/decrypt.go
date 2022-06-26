package dash

import (
   "github.com/edgeware/mp4ff/mp4"
   "io"
)

type decrypter struct {
   *mp4.SinfBox
   w io.Writer
}

func new_decrypter(w io.Writer) decrypter {
   return decrypter{w: w}
}

func (d *decrypter) init(r io.Reader) error {
   file, err := mp4.DecodeFile(r)
   if err != nil {
      return err
   }
   for _, trak := range file.Init.Moov.Traks {
      for _, child := range trak.Mdia.Minf.Stbl.Stsd.Children {
         switch child.Type() {
         case "encv":
            d.SinfBox, err = child.(*mp4.VisualSampleEntryBox).RemoveEncryption()
         case "enca":
            d.SinfBox, err = child.(*mp4.AudioSampleEntryBox).RemoveEncryption()
         }
         if err != nil {
            return err
         }
      }
   }
   file.Init.Moov.RemovePsshs()
   return file.Init.Encode(d.w)
}

func (d decrypter) segment(r io.Reader, key []byte) error {
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
            tenc := d.SinfBox.Schi.Tenc
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
               switch d.SinfBox.Schm.SchemeType {
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
      err := seg.Encode(d.w)
      if err != nil {
         return err
      }
   }
   return nil
}
