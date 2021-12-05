package m3u

import (
   "fmt"
   "os"
   "testing"
)

var tests = []string{
   "2.m3u8",
   "HLSPlaylist.m3u8",
   "HLS_540.m3u8",
   "pc_hd_abr_v2_hls_master.m3u8",
   "vf_p0b7tn99_7a17a782-83d5-43f2-8387-d3008bc6f2c1-audio_eng=96000-video=1604000.m3u8",
}

func TestPlaylist(t *testing.T) {
   for _, test := range tests {
      fmt.Println(test + ":")
      buf, err := os.ReadFile(test)
      if err != nil {
         t.Fatal(err)
      }
      for _, form := range Unmarshal(buf, "dir/") {
         fmt.Printf("%+v\n", form)
      }
      fmt.Println()
   }
}
