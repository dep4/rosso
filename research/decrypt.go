package dash

import (
   "crypto/aes"
   "crypto/cipher"
   "github.com/edgeware/mp4ff/mp4"
   "io"
)

func Decrypt(w io.Writer, r io.Reader, key []byte) error {
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
            sinf := file.Moov.GetSinf(traf.Tfhd.TrackID)
            if sinf == nil {
               continue
            }
            iv := sinf.Schi.Tenc.DefaultConstantIV
            for i, sample := range samples {
               var sub []mp4.SubSamplePattern
               if len(traf.Senc.SubSamples) > i {
                  sub = traf.Senc.SubSamples[i]
               }
               dec, err := DecryptSampleCenc(sample.Data, key, iv, sub)
               if err != nil {
                  return err
               }
               copy(sample.Data, dec)
            }
            traf.RemoveEncryptionBoxes()
         }
      }
      err := seg.Encode(w)
      if err != nil {
         return err
      }
   }
   return nil
}

func min(a, b uint32) uint32 {
   if a < b {
      return a
   }
   return b
}

func DecryptBytes(data []byte, key []byte, iv []byte) ([]byte, error) {
   block, err := aes.NewCipher(key)
   if err != nil {
      return nil, err
   }
   cipher.NewCBCDecrypter(block, iv).CryptBlocks(data, data)
   return data, nil
}

func DecryptSampleCenc(input []byte, key []byte, iv []byte, subSamplePatterns []mp4.SubSamplePattern) ([]byte, error) {
   var decSample []byte
   for _, ss := range subSamplePatterns {
      input := ss.BytesOfClearData
      rem_bytes := ss.BytesOfProtectedData
      for rem_bytes > 0 {
         // aomediacodec.github.io/av1-isobmff
         // crypt_byte_block = 1 and skip_byte_block = 9
         if rem_bytes < 16*1 {
            break
         }
         //cryptOut, err := DecryptBytesCTR(input[pos:pos+rem_bytes], key, iv)
         cryptOut, err := DecryptBytes(input, key, iv)
         if err != nil {
            return nil, err
         }
         decSample = append(decSample, sample[pos:pos+input]...)
         pos += input
         decSample = append(decSample, cryptOut...)
         pos += rem_bytes
         rem_bytes -= min(16*9, rem_bytes)
      }
   }
   return decSample, nil
}
