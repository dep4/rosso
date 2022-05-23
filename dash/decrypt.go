package dash

import (
   "github.com/edgeware/mp4ff/mp4"
   "io"
   "strconv"
   "strings"
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
            for i, sample := range samples {
               var iv []byte
               // this needs its own line so that the bytes are copied
               iv = append(iv, traf.Senc.IVs[i]...)
               iv = append(iv, 0, 0, 0, 0, 0, 0, 0, 0)
               var sub []mp4.SubSamplePattern
               if len(traf.Senc.SubSamples) > i {
                  sub = traf.Senc.SubSamples[i]
               }
               dec, err := mp4.DecryptSampleCenc(sample.Data, key, iv, sub)
               if err != nil {
                  return err
               }
               copy(sample.Data, dec)
            }
            // required for playback
            traf.RemoveEncryptionBoxes()
         }
         // fast start
         frag.Moof.RemovePsshs()
      }
      // fix jerk between fragments
      seg.Sidx = nil
      err := seg.Encode(w)
      if err != nil {
         return err
      }
   }
   return nil
}

func (r Represent) id(in string) string {
   return strings.Replace(in, "$RepresentationID$", r.ID, 1)
}

func (s Segment) number(in string) string {
   return strings.Replace(in, "$Number$", strconv.Itoa(s.T), 1)
}

func (s Segment) time(in string) string {
   return strings.Replace(in, "$Time$", strconv.Itoa(s.T), 1)
}
