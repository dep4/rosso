package m3u8

import (
   "fmt"
   "github.com/89z/format"
   "io"
   "net/http"
   "net/url"
   "os"
   "path"
   "path/filepath"
   "strings"
   "time"
)

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
