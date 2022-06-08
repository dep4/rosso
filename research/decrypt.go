package decrypt

import (
   "bytes"
   "crypto/aes"
   "crypto/cipher"
   "fmt"
   "github.com/edgeware/mp4ff/mp4"
   "io"
)

func decryptAndWriteSegments(inMp4 *mp4.File, tracks []trackInfo, key []byte, ofh io.Writer) error {
   var outNr uint32 = 1
   for _, seg := range inMp4.Segments {
      for _, frag := range seg.Fragments {
         err := decryptFragment(inMp4, frag, tracks, key)
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
      for _, child := range stsd.Children {
         switch child.Type() {
         case "encv":
            encv = child.(*mp4.VisualSampleEntryBox)
            sinf, err := encv.RemoveEncryption()
            if err != nil {
               return err
            }
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
            tracks = append(tracks, trackInfo{
               trackID: trackID,
               sinf:    sinf,
            })
         default:
            continue
         }
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
   err = decryptAndWriteSegments(inMp4, tracks, key, w)
   if err != nil {
      return err
   }
   return nil
}

func decryptFragment(inMp4 *mp4.File, frag *mp4.Fragment, tracks []trackInfo, key []byte) error {
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
      var iv []byte
      if sinf := inMp4.Moov.GetSinf(traf.Tfhd.TrackID); sinf != nil {
         iv = sinf.Schi.Tenc.DefaultConstantIV
      }
      if iv == nil {
         continue
      }
      samples, err := frag.GetFullSamples(ti.trex)
      if err != nil {
         return err
      }
      err = decryptSamplesInPlace(samples, key, iv, traf.Senc)
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
   _ = moof.RemovePsshs()
   return nil
}

func decryptSamplesInPlace(samples []mp4.FullSample, key, iv []byte, senc *mp4.SencBox) error {
   for i := range samples {
      encSample := samples[i].Data
      var subSamplePatterns []mp4.SubSamplePattern
      if len(senc.SubSamples) != 0 {
         subSamplePatterns = senc.SubSamples[i]
      }
      decryptedSample, err := DecryptSample(encSample, key, iv, subSamplePatterns)
      if err != nil {
         return err
      }
      _ = copy(samples[i].Data, decryptedSample)
   }
   return nil
}

func DecryptSample(input, key, iv []byte, subsamples []mp4.SubSamplePattern) ([]byte, error) {
   var (
      output []byte
      pos uint32
   )
   for _, subsample := range subsamples {
      output = append(output, input[pos:][:subsample.BytesOfClearData]...)
      pos += uint32(subsample.BytesOfClearData)
      rem_bytes := subsample.BytesOfProtectedData
      for rem_bytes > 0 {
         // aomediacodec.github.io/av1-isobmff
         // crypt_byte_block = 1 and skip_byte_block = 9
         if rem_bytes < 16*1 {
            output = append(output, input[pos:]...)
            break
         }
         data, err := DecryptBytes(input[pos:][:16*1], key, iv)
         if err != nil {
            return nil, err
         }
         output = append(output, data...)
         // crypt
         pos += 16*1
         rem_bytes -= 16*1
         // skip
         pos += min(16*9, rem_bytes)
         rem_bytes -= min(16*9, rem_bytes)
      }
      pos += subsample.BytesOfProtectedData
   }
   return output, nil
}

func DecryptBytes(data []byte, key []byte, iv []byte) ([]byte, error) {
   block, err := aes.NewCipher(key)
   if err != nil {
      return nil, err
   }
   cipher.NewCBCDecrypter(block, iv).CryptBlocks(data, data)
   return data, nil
}

func min(a, b uint32) uint32 {
   if a < b {
      return a
   }
   return b
}

type trackInfo struct {
	trackID uint32
	sinf    *mp4.SinfBox
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

