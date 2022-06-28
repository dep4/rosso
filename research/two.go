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
   moov := file.Init.Moov
   for _, trak := range moov.Traks {
      stsd := trak.Mdia.Minf.Stbl.Stsd
      for _, child := range stsd.Children {
         switch child.Type() {
         case "encv":
            encv := child.(*mp4.VisualSampleEntryBox)
            d.SinfBox, err = encv.RemoveEncryption()
            if err != nil {
               return err
            }
         case "enca":
            enca := child.(*mp4.AudioSampleEntryBox)
            d.SinfBox, err = enca.RemoveEncryption()
            if err != nil {
               return err
            }
         }
      }
   }
   moov.RemovePsshs()
   return file.Init.Encode(d.w)
}

func (d decrypter) segment(r io.Reader, key []byte) error {
   file, err := mp4.DecodeFile(r)
   if err != nil {
      return err
   }
   for _, seg := range file.Segments {
      for _, frag := range seg.Fragments {
         moof := frag.Moof
         var nrBytesRemoved uint64 = 0
         for _, traf := range moof.Trafs {
            schemeType := d.SinfBox.Schm.SchemeType
            tenc := d.SinfBox.Schi.Tenc
            samples, err := frag.GetFullSamples(nil)
            if err != nil {
               return err
            }
            for i := range samples {
               encSample := samples[i].Data
               var iv []byte
               if len(traf.Senc.IVs) == len(samples) {
                  if len(traf.Senc.IVs[i]) == 8 {
                     iv = make([]byte, 0, 16)
                     iv = append(iv, traf.Senc.IVs[i]...)
                     iv = append(iv, []byte{0, 0, 0, 0, 0, 0, 0, 0}...)
                  } else if len(traf.Senc.IVs) == len(samples) {
                     iv = traf.Senc.IVs[i]
                  }
               } else if tenc.DefaultConstantIV != nil {
                  iv = tenc.DefaultConstantIV
               }
               var subSamplePatterns []mp4.SubSamplePattern
               if len(traf.Senc.SubSamples) != 0 {
                  subSamplePatterns = traf.Senc.SubSamples[i]
               }
               switch schemeType {
               case "cenc":
                  err := mp4.DecryptSampleCenc(encSample, key, iv, subSamplePatterns)
                  if err != nil {
                     return err
                  }
               case "cbcs":
                  err := mp4.DecryptSampleCbcs(encSample, key, iv, subSamplePatterns, tenc)
                  if err != nil {
                     return err
                  }
               }
            }
            nrBytesRemoved += traf.RemoveEncryptionBoxes()
         }
         _, psshBytesRemoved := moof.RemovePsshs()
         nrBytesRemoved += psshBytesRemoved
         for _, traf := range moof.Trafs {
            for _, trun := range traf.Truns {
               trun.DataOffset -= int32(nrBytesRemoved)
            }
         }
      }
      seg.Sidx = nil // drop sidx inside segment, since not modified properly
      err := seg.Encode(d.w)
      if err != nil {
         return err
      }
   }
   return nil
}
