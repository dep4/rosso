package m3u8

import (
   "errors"
   "fmt"
   "io"
   "net/url"
   "os"
   "path/filepath"
   "regexp"
   "path"
   "net/http"
   "github.com/89z/format"
   "strconv"
   "strings"
   "sync"
   "time"
)

// NewTask returns a Task instance
func NewTask(output string, url string) (*Downloader, error) {
   result, err := FromURL(url)
   if err != nil {
   return nil, err
   }
   var folder string
   // If no output folder specified, use current directory
   if output == "" {
   current, err := CurrentDir()
   if err != nil {
   return nil, err
   }
   folder = filepath.Join(current, output)
   } else {
   folder = output
   }
   if err := os.MkdirAll(folder, os.ModePerm); err != nil {
   return nil, fmt.Errorf("create storage folder failed: %s", err.Error())
   }
   tsFolder := filepath.Join(folder, tsFolderName)
   if err := os.MkdirAll(tsFolder, os.ModePerm); err != nil {
   return nil, fmt.Errorf("create ts folder '[%s]' failed: %s", tsFolder, err.Error())
   }
   d := &Downloader{
   folder:   folder,
   tsFolder: tsFolder,
   result:   result,
   }
   d.segLen = len(result.M3u8.Segments)
   d.queue = genSlice(d.segLen)
   return d, nil
}

// Start runs downloader
func (d *Downloader) Start(concurrency int) error {
   var wg sync.WaitGroup
   // struct{} zero size
   limitChan := make(chan struct{}, concurrency)
   for {
      tsIdx, end, err := d.next()
      if err != nil {
         if end {
            break
         }
         continue
      }
      wg.Add(1)
      go func(idx int) {
         defer wg.Done()
         if err := d.download(idx); err != nil {
            // Back into the queue, retry request
            fmt.Printf("[failed] %s\n", err.Error())
            if err := d.back(idx); err != nil {
               fmt.Println(err)
            }
         }
         <-limitChan
      }(tsIdx)
      limitChan <- struct{}{}
   }
   wg.Wait()
   return nil
}

////////////////////////////////////////////////////////////////////////////////

const (
	tsExt            = ".ts"
	tsFolderName     = "ts"
	mergeTSFilename  = "main.ts"
	tsTempFileSuffix = "_tmp"
	progressWidth    = 40
)

type Downloader struct {
	lock     sync.Mutex
	queue    []int
	folder   string
	tsFolder string
	finish   int32
	segLen   int

	result *Result
}

func (d *Downloader) next() (segIndex int, end bool, err error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if len(d.queue) == 0 {
		err = fmt.Errorf("queue empty")
		if d.finish == int32(d.segLen) {
			end = true
			return
		}
		// Some segment indexes are still running.
		end = false
		return
	}
	segIndex = d.queue[0]
	d.queue = d.queue[1:]
	return
}

func (d *Downloader) back(segIndex int) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	if sf := d.result.M3u8.Segments[segIndex]; sf == nil {
		return fmt.Errorf("invalid segment index: %d", segIndex)
	}
	d.queue = append(d.queue, segIndex)
	return nil
}

func (d *Downloader) tsURL(segIndex int) string {
	seg := d.result.M3u8.Segments[segIndex]
	return ResolveURL(d.result.URL, seg.URI)
}

func tsFilename(ts int) string {
	return strconv.Itoa(ts) + tsExt
}

func genSlice(len int) []int {
	s := make([]int, 0)
	for i := 0; i < len; i++ {
		s = append(s, i)
	}
	return s
}

type Result struct {
	URL  *url.URL
	M3u8 *M3u8
	Keys map[int]string
}

func CurrentDir(joinPath ...string) (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}
	p := strings.Replace(dir, "\\", "/", -1)
	whole := filepath.Join(joinPath...)
	whole = filepath.Join(p, whole)
	return whole, nil
}

func ResolveURL(u *url.URL, p string) string {
	if strings.HasPrefix(p, "https://") || strings.HasPrefix(p, "http://") {
		return p
	}
	var baseURL string
	if strings.Index(p, "/") == 0 {
		baseURL = u.Scheme + "://" + u.Host
	} else {
		tU := u.String()
		baseURL = tU[0:strings.LastIndex(tU, "/")]
	}
	return baseURL + path.Join("/", p)
}

func DrawProgressBar(prefix string, proportion float32, width int, suffix ...string) {
	pos := int(proportion * float32(width))
	s := fmt.Sprintf("[%s] %s%*s %6.2f%% %s",
		prefix, strings.Repeat("■", pos), width-pos, "", proportion*100, strings.Join(suffix, ""))
	fmt.Print("\r" + s)
}

func Get(url string) (io.ReadCloser, error) {
   c := http.Client{
   Timeout: time.Duration(60) * time.Second,
   }
   req, err := http.NewRequest("GET", url, nil)
   if err != nil {
   return nil, err
   }
   req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36")
   format.LogLevel.Dump(1, req)
   resp, err := c.Do(req)
   if err != nil {
   return nil, err
   }
   if resp.StatusCode != 200 {
   return nil, fmt.Errorf("http error: status code %d", resp.StatusCode)
   }
   return resp.Body, nil
}

type (
	PlaylistType string
	CryptMethod  string
)

const (
   PlaylistTypeEvent PlaylistType = "EVENT"
   PlaylistTypeVOD   PlaylistType = "VOD"
)

// regex pattern for extracting `key=value` parameters from a line
var linePattern = regexp.MustCompile(`([a-zA-Z-]+)=("[^"]+"|[^",]+)`)

type M3u8 struct {
	Version        int8   // EXT-X-VERSION:version
	MediaSequence  uint64 // Default 0, #EXT-X-MEDIA-SEQUENCE:sequence
	Segments       []*Segment
	MasterPlaylist []*MasterPlaylist
	Keys           map[int]*Key
	EndList        bool         // #EXT-X-ENDLIST
	PlaylistType   PlaylistType // VOD or EVENT
	TargetDuration float64      // #EXT-X-TARGETDURATION:duration
}

type Segment struct {
	URI      string
	KeyIndex int
	Title    string  // #EXTINF: duration,<title>
	Duration float32 // #EXTINF: duration,<title>
	Length   uint64  // #EXT-X-BYTERANGE: length[@offset]
	Offset   uint64  // #EXT-X-BYTERANGE: length[@offset]
}

// #EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=240000,RESOLUTION=416x234,CODECS="avc1.42e00a,mp4a.40.2"
type MasterPlaylist struct {
	URI        string
	BandWidth  uint32
	Resolution string
	Codecs     string
	ProgramID  uint32
}

func parseMasterPlaylist(line string) (*MasterPlaylist, error) {
	params := parseLineParameters(line)
	if len(params) == 0 {
		return nil, errors.New("empty parameter")
	}
	mp := new(MasterPlaylist)
	for k, v := range params {
		switch {
		case k == "BANDWIDTH":
			v, err := strconv.ParseUint(v, 10, 32)
			if err != nil {
				return nil, err
			}
			mp.BandWidth = uint32(v)
		case k == "RESOLUTION":
			mp.Resolution = v
		case k == "PROGRAM-ID":
			v, err := strconv.ParseUint(v, 10, 32)
			if err != nil {
				return nil, err
			}
			mp.ProgramID = uint32(v)
		case k == "CODECS":
			mp.Codecs = v
		}
	}
	return mp, nil
}

// parseLineParameters extra parameters in string `line`
func parseLineParameters(line string) map[string]string {
	r := linePattern.FindAllStringSubmatch(line, -1)
	params := make(map[string]string)
	for _, arr := range r {
		params[arr[1]] = strings.Trim(arr[2], "\"")
	}
	return params
}
