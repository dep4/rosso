package decrypt

import (
   "bytes"
   "fmt"
   "github.com/edgeware/mp4ff/mp4"
   "io"
)

func decryptInit(r io.Reader, w io.Writer) error {
   inMp4, err := mp4.DecodeFile(r)
   if err != nil {
      return err
   }
   if !inMp4.IsFragmented() {
      return fmt.Errorf("file not fragmented. Not supported")
   }
   moov := inMp4.Init.Moov
   for _, trak := range moov.Traks {
      stsd := trak.Mdia.Minf.Stbl.Stsd
      for _, child := range stsd.Children {
         switch child.Type() {
         case "encv":
            encv := child.(*mp4.VisualSampleEntryBox)
            _, err := encv.RemoveEncryption()
            if err != nil {
               return err
            }
         case "enca":
            enca := child.(*mp4.AudioSampleEntryBox)
            _, err := enca.RemoveEncryption()
            if err != nil {
               return err
            }
         }
      }
   }
   psshs := moov.RemovePsshs()
   for _, pssh := range psshs {
      psshInfo := bytes.Buffer{}
      err = pssh.Info(&psshInfo, "", "", "  ")
      if err != nil {
         return err
      }
   }
   return inMp4.Init.Encode(w)
}

func decryptSegment(r io.Reader, key []byte, w io.Writer) error {
   inMp4, err := mp4.DecodeFile(r)
   if err != nil {
      return err
   }
   if !inMp4.IsFragmented() {
      return fmt.Errorf("file not fragmented. Not supported")
   }
   var tracks []trackInfo
   for _, seg := range inMp4.Segments {
      for _, frag := range seg.Fragments {
         moof := frag.Moof
         var nrBytesRemoved uint64
         for _, traf := range moof.Trafs {
            hasSenc, isParsed := traf.ContainsSencBox()
            if !hasSenc {
               return fmt.Errorf("no senc box in traf")
            }
            var ti trackInfo
            for _, track := range tracks {
               if track.trackID == traf.Tfhd.TrackID {
                  ti = track
               }
            }
            if !isParsed {
               defaultIVSize := ti.sinf.Schi.Tenc.DefaultPerSampleIVSize
               err := traf.ParseReadSenc(defaultIVSize, moof.StartPos)
               if err != nil {
                  return fmt.Errorf("parseReadSenc: %w", err)
               }
            }
            samples, err := frag.GetFullSamples(ti.trex)
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
               var subSamplePatterns []mp4.SubSamplePattern
               if len(traf.Senc.SubSamples) != 0 {
                  subSamplePatterns = traf.Senc.SubSamples[i]
               }
               decryptedSample, err := mp4.DecryptSampleCenc(encSample, key, iv, subSamplePatterns)
               if err != nil {
                  return err
               }
               copy(samples[i].Data, decryptedSample)
            }
            nrBytesRemoved += traf.RemoveEncryptionBoxes()
         }
         for _, traf := range moof.Trafs {
            for _, trun := range traf.Truns {
               trun.DataOffset -= int32(nrBytesRemoved)
            }
         }
         moof.RemovePsshs()
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

/////////////////////////////////////////////////////////////////////

type trackInfo struct {
   sinf    *mp4.SinfBox
   trackID uint32
   trex    *mp4.TrexBox
}
