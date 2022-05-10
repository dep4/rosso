package dash

import (
   "fmt"
   "github.com/edgeware/mp4ff/mp4"
   "io"
)

func decryptMP4withCenc(r io.Reader, key []byte, w io.Writer) error {
   inMp4, err := mp4.DecodeFile(r)
   if err != nil {
      return err
   }
   var tracks []trackInfo
   for _, trak := range inMp4.Init.Moov.Traks {
      for _, child := range trak.Mdia.Minf.Stbl.Stsd.Children {
         var box interface {
            RemoveEncryption() (*mp4.SinfBox, error)
         }
         switch child.Type() {
         case "encv":
            box = child.(*mp4.VisualSampleEntryBox)
         case "enca":
            box = child.(*mp4.AudioSampleEntryBox)
         }
         sinf, err := box.RemoveEncryption()
         if err != nil {
            return err
         }
         tracks = append(tracks, trackInfo{
            sinf:    sinf,
            trackID: trak.Tkhd.TrackID,
         })
      }
   }
   for _, trex := range inMp4.Init.Moov.Mvex.Trexs {
      for i := range tracks {
         if tracks[i].trackID == trex.TrackID {
            tracks[i].trex = trex
            break
         }
      }
   }
   if err := inMp4.Init.Encode(w); err != nil {
      return err
   }
   return decryptAndWriteSegments(inMp4.Segments, tracks, key, w)
}

func decryptSamplesInPlace(samples []mp4.FullSample, key []byte, senc *mp4.SencBox) error {
   for i, sam := range samples {
      var (
         iv []byte
         subSamplePatterns []mp4.SubSamplePattern
      )
      if len(senc.IVs[i]) == 8 {
         iv = append(iv, senc.IVs[i]...)
         iv = append(iv, 0, 0, 0, 0, 0, 0, 0, 0)
      } else {
         iv = senc.IVs[i]
      }
      if len(senc.SubSamples) != 0 {
         subSamplePatterns = senc.SubSamples[i]
      }
      dec, err := mp4.DecryptSampleCenc(sam.Data, key, iv, subSamplePatterns)
      if err != nil {
         return err
      }
      copy(samples[i].Data, dec)
   }
   return nil
}

func decryptAndWriteSegments(segs []*mp4.MediaSegment, tracks []trackInfo, key []byte, ofh io.Writer) error {
   for _, seg := range segs {
      for _, frag := range seg.Fragments {
         err := decryptFragment(frag, tracks, key)
         if err != nil {
            return err
         }
      }
      if seg.Sidx != nil {
         seg.Sidx = nil // drop sidx inside segment, since not modified properly
      }
      err := seg.Encode(ofh)
      if err != nil {
         return err
      }
   }
   return nil
}

type trackInfo struct {
   trackID uint32
   sinf    *mp4.SinfBox
   trex    *mp4.TrexBox
}

func decryptFragment(frag *mp4.Fragment, tracks []trackInfo, key []byte) error {
   var nrBytesRemoved uint64
   for _, traf := range frag.Moof.Trafs {
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
         err := traf.ParseReadSenc(defaultIVSize, frag.Moof.StartPos)
         if err != nil {
            return fmt.Errorf("parseReadSenc: %w", err)
         }
      }
      samples, err := frag.GetFullSamples(ti.trex)
      if err != nil {
         return err
      }
      if err := decryptSamplesInPlace(samples, key, traf.Senc); err != nil {
         return err
      }
      nrBytesRemoved += traf.RemoveEncryptionBoxes()
   }
   return nil
}
