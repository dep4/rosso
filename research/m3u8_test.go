package m3u8

import (
	"flag"
	"fmt"
	"os"
)

/*
https://cbsios-vh.akamaihd.net/i/temp_hd_gallery_video/CBS_Production_Outlet_VMS/video_robot/CBS_Production_Entertainment/2021/10/18/1963091011554/NICKELODEON_STARTREKPRODIGY_104_HD_985058_,2228,4628,3128,1628,848,503,000.mp4.csmil/index_1_av.m3u8?null=0&id=AgBItRcmFy81SreI62E6ac+Xz62gjgxyv1YRgnd02LQKlqbMvUT+1mZJ22FHNdI832cTQwH6Xz+CEg%3d%3d&hdntl=exp=1642912311~acl=%2fi%2ftemp_hd_gallery_video%2fCBS_Production_Outlet_VMS%2fvideo_robot%2fCBS_Production_Entertainment%2f2021%2f10%2f18%2f1963091011554%2fNICKELODEON_STARTREKPRODIGY_104_HD_985058_*~data=hdntl~hmac=d1e1c2257b1e3c98e31225cbed1605a313eb6e18286287dd8ccbcca785c5166f
*/

var (
	url      string
	output   string
	chanSize int
)

func init() {
	flag.StringVar(&url, "u", "", "M3U8 URL, required")
	flag.IntVar(&chanSize, "c", 25, "Maximum number of occurrences")
	flag.StringVar(&output, "o", "", "Output folder, required")
}

func main() {
   flag.Parse()
   defer func() {
   if r := recover(); r != nil {
   fmt.Println("[error]", r)
   os.Exit(-1)
   }
   }()
   if url == "" {
   panicParameter("u")
   }
   if output == "" {
   panicParameter("o")
   }
   if chanSize <= 0 {
   panic("parameter 'c' must be greater than 0")
   }
   downloader, err := NewTask(output, url)
   if err != nil {
   panic(err)
   }
   if err := downloader.Start(chanSize); err != nil {
   panic(err)
   }
   fmt.Println("Done!")
}

func panicParameter(name string) {
	panic("parameter '" + name + "' is required")
}
