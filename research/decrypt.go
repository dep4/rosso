package dash

import (
   "github.com/edgeware/mp4ff/mp4"
   "io"
)

type track_info struct {
   id uint32
   sinf *mp4.SinfBox
}

type decrypter struct {
   //tracks []track_info
   tracks map[uint32]*mp4.SinfBox
   w io.Writer
}

func new_decrypter(w io.Writer) decrypter {
   var d decrypter
   d.tracks = make(map[uint32]*mp4.SinfBox)
   d.w = w
   return d
}

func (d *decrypter) init(r io.Reader) error {
   file, err := mp4.DecodeFile(r)
   if err != nil {
      return err
   }
   moov := file.Init.Moov
   for _, trak := range moov.Traks {
      trackID := trak.Tkhd.TrackID
      stsd := trak.Mdia.Minf.Stbl.Stsd
      for _, child := range stsd.Children {
         switch child.Type() {
         case "encv":
            sinf, err := child.(*mp4.VisualSampleEntryBox).RemoveEncryption()
            if err != nil {
               return err
            }
            d.tracks[trackID] = sinf
            /*
            d.tracks = append(d.tracks, track_info{
               id: trackID,
               sinf:    sinf,
            })
            */
         case "enca":
            sinf, err := child.(*mp4.AudioSampleEntryBox).RemoveEncryption()
            if err != nil {
               return err
            }
            d.tracks[trackID] = sinf
            /*
            d.tracks = append(d.tracks, track_info{
               id: trackID,
               sinf:    sinf,
            })
            */
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
         var nrBytesRemoved uint64
         for _, traf := range frag.Moof.Trafs {
            sinf := d.tracks[traf.Tfhd.TrackID]
            if sinf == nil {
               continue
            }
            tenc := sinf.Schi.Tenc
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
               var sub []mp4.SubSamplePattern
               if len(traf.Senc.SubSamples) != 0 {
                  sub = traf.Senc.SubSamples[i]
               }
               switch sinf.Schm.SchemeType {
               case "cenc":
                  err = mp4.DecryptSampleCenc(encSample, key, iv, sub)
               case "cbcs":
                  err = mp4.DecryptSampleCbcs(encSample, key, iv, sub, tenc)
               }
               if err != nil {
                  return err
               }
            }
            nrBytesRemoved += traf.RemoveEncryptionBoxes()
         }
         _, psshBytesRemoved := frag.Moof.RemovePsshs()
         nrBytesRemoved += psshBytesRemoved
         for _, traf := range frag.Moof.Trafs {
            for _, trun := range traf.Truns {
               trun.DataOffset -= int32(nrBytesRemoved)
            }
         }
      }
      seg.Sidx = nil
      err := seg.Encode(d.w)
      if err != nil {
         return err
      }
   }
   return nil
}
