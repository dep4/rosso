package dash

import (
   "bytes"
   "fmt"
   "github.com/edgeware/mp4ff/mp4"
   "io"
)

func decryptMP4withCenc(r io.Reader, key []byte, w io.Writer) error {
   inMp4, err := mp4.DecodeFile(r)
   if err != nil {
      return err
   }
   if !inMp4.IsFragmented() {
      return fmt.Errorf("file not fragmented. Not supported")
   }
   tracks := make([]trackInfo, 0, len(inMp4.Init.Moov.Traks))
   moov := inMp4.Init.Moov
   for _, trak := range moov.Traks {
      trackID := trak.Tkhd.TrackID
      stsd := trak.Mdia.Minf.Stbl.Stsd
      var encv *mp4.VisualSampleEntryBox
      var enca *mp4.AudioSampleEntryBox
      var schemeType string
      for _, child := range stsd.Children {
         switch child.Type() {
         case "encv":
            encv = child.(*mp4.VisualSampleEntryBox)
            sinf, err := encv.RemoveEncryption()
            if err != nil {
               return err
            }
            schemeType = sinf.Schm.SchemeType
            tracks = append(tracks, trackInfo{
               trackID: trackID,
               sinf:    sinf,
            })
         case "enca":
            enca = child.(*mp4.AudioSampleEntryBox)
            sinf, err := enca.RemoveEncryption()
            if err != nil {
               return err
            }
            schemeType = sinf.Schm.SchemeType
            tracks = append(tracks, trackInfo{
               trackID: trackID,
               sinf:    sinf,
            })
         default:
            continue
         }
      }
      if schemeType != "" && schemeType != "cenc" && schemeType != "cbcs" {
         return fmt.Errorf("scheme type %s not supported", schemeType)
      }
      if schemeType == "" {
         tracks = append(tracks, trackInfo{
            trackID: trackID,
            sinf:    nil,
         })
      }
   }
   for _, trex := range moov.Mvex.Trexs {
      for i := range tracks {
         if tracks[i].trackID == trex.TrackID {
            tracks[i].trex = trex
            break
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
   err = inMp4.Init.Encode(w)
   if err != nil {
      return err
   }
   var outNr uint32 = 1
   for _, seg := range inMp4.Segments {
      for _, frag := range seg.Fragments {
         moof := frag.Moof
         var nrBytesRemoved uint64
         for _, traf := range moof.Trafs {
            var ti trackInfo
            for _, track := range tracks {
               if track.trackID == traf.Tfhd.TrackID {
                  ti = track
               }
            }
            if ti.sinf != nil {
               schemeType := ti.sinf.Schm.SchemeType
               if schemeType != "cenc" && schemeType != "cbcs" {
                  return fmt.Errorf("scheme type %s not supported", schemeType)
               }
               hasSenc, isParsed := traf.ContainsSencBox()
               if !hasSenc {
                  return fmt.Errorf("no senc box in traf")
               }
               if !isParsed {
                  defaultPerSampleIVSize := ti.sinf.Schi.Tenc.DefaultPerSampleIVSize
                  err := traf.ParseReadSenc(defaultPerSampleIVSize, moof.StartPos)
                  if err != nil {
                     return fmt.Errorf("parseReadSenc: %w", err)
                  }
               }
               tenc := ti.sinf.Schi.Tenc
               samples, err := frag.GetFullSamples(ti.trex)
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
                  if len(iv) == 0 {
                     return fmt.Errorf("iv has length 0")
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
         }
         _, psshBytesRemoved := moof.RemovePsshs()
         nrBytesRemoved += psshBytesRemoved
         for _, traf := range moof.Trafs {
            for _, trun := range traf.Truns {
               trun.DataOffset -= int32(nrBytesRemoved)
            }
         }
         outNr++
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

type trackInfo struct {
   sinf    *mp4.SinfBox
   trackID uint32
   trex    *mp4.TrexBox
}
