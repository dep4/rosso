package decrypt

import (
   "bytes"
   "crypto/aes"
   "crypto/cipher"
   "fmt"
   "github.com/edgeware/mp4ff/mp4"
   "io"
)

// aomediacodec.github.io/av1-isobmff
// crypt_byte_block = 1 and skip_byte_block = 9
func DecryptSample(input, key, iv []byte, subsamples []mp4.SubSamplePattern) ([]byte, error) {
   block, err := aes.NewCipher(key)
   if err != nil {
      return nil, err
   }
   var low uint32
   for _, subsample := range subsamples {
      for {
         high := low + 16*1
         if int(high) >= len(input) {
            break
         }
         data := input[low:high]
         cipher.NewCBCDecrypter(block, iv).CryptBlocks(data, data)
         // crypt
         low += 16*1
         subsample.BytesOfProtectedData -= 16*1
         // skip
         low += 16*9
         subsample.BytesOfClearData -= 16*9
      }
      low += subsample.BytesOfProtectedData
      low += uint32(subsample.BytesOfClearData)
   }
   return input, nil
}

func decryptFragment(inMp4 *mp4.File, frag *mp4.Fragment, tracks []trackInfo, key []byte) error {
   moof := frag.Moof
   var nrBytesRemoved uint64
   for _, traf := range moof.Trafs {
      hasSenc, isParsed := traf.ContainsSencBox()
      if !hasSenc {
         continue
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
      iv := ti.sinf.Schi.Tenc.DefaultConstantIV
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
   var tracks []trackInfo
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

