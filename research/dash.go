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
   if !inMp4.IsFragmented() {
      return fmt.Errorf("file not fragmented. Not supported")
   }
   var tracks []trackInfo
   return decryptAndWriteSegments(inMp4.Segments, tracks, key, w)
}

func decryptAndWriteSegments(segs []*mp4.MediaSegment, tracks []trackInfo, key []byte, ofh io.Writer) error {
   var outNr uint32 = 1
   for _, seg := range segs {
      for _, frag := range seg.Fragments {
         err := decryptFragment(frag, tracks, key)
         if err != nil {
            return err
         }
         outNr++
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
   sinf    *mp4.SinfBox
   trackID uint32
   trex    *mp4.TrexBox
}

func findTrackInfo(tracks []trackInfo, trackID uint32) trackInfo {
   for _, ti := range tracks {
      if ti.trackID == trackID {
         return ti
      }
   }
   return trackInfo{}
}

func decryptFragment(frag *mp4.Fragment, tracks []trackInfo, key []byte) error {
   moof := frag.Moof
   var nrBytesRemoved uint64 = 0
   for _, traf := range moof.Trafs {
      hasSenc, isParsed := traf.ContainsSencBox()
      if !hasSenc {
         return fmt.Errorf("no senc box in traf")
      }
      ti := findTrackInfo(tracks, traf.Tfhd.TrackID)
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
      err = decryptSamplesInPlace(samples, key, traf.Senc)
      if err != nil {
         return err
      }
      nrBytesRemoved += traf.RemoveEncryptionBoxes()
   }
   for _, traf := range moof.Trafs {
      for _, trun := range traf.Truns {
         trun.DataOffset -= int32(nrBytesRemoved)
      }
   }
   moof.RemovePsshs()
   return nil
}

func decryptSamplesInPlace(samples []mp4.FullSample, key []byte, senc *mp4.SencBox) error {
   for i := range samples {
      encSample := samples[i].Data
      var iv []byte
      if len(senc.IVs[i]) == 8 {
         iv = make([]byte, 0, 16)
         iv = append(iv, senc.IVs[i]...)
         iv = append(iv, []byte{0, 0, 0, 0, 0, 0, 0, 0}...)
      } else {
         iv = senc.IVs[i]
      }
      var subSamplePatterns []mp4.SubSamplePattern
      if len(senc.SubSamples) != 0 {
         subSamplePatterns = senc.SubSamples[i]
      }
      decryptedSample, err := mp4.DecryptSampleCenc(encSample, key, iv, subSamplePatterns)
      if err != nil {
         return err
      }
      copy(samples[i].Data, decryptedSample)
   }
   return nil
}
