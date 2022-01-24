package m3u8

import (
   "bufio"
   "bytes"
   "crypto/aes"
   "crypto/cipher"
   "errors"
   "fmt"
   "io"
   "io/ioutil"
   "net/url"
   "os"
   "path/filepath"
   "strconv"
   "strings"
   "sync/atomic"
)

func pkcs5UnPadding(origData []byte) []byte {
   length := len(origData)
   unPadding := int(origData[length-1])
   return origData[:(length - unPadding)]
}

func AES128Decrypt(crypted, key, iv []byte) ([]byte, error) {
   block, err := aes.NewCipher(key)
   if err != nil {
      return nil, err
   }
   blockSize := block.BlockSize()
   if len(iv) == 0 {
      iv = key
   }
   blockMode := cipher.NewCBCDecrypter(block, iv[:blockSize])
   origData := make([]byte, len(crypted))
   blockMode.CryptBlocks(origData, crypted)
   origData = pkcs5UnPadding(origData)
   return origData, nil
}

func (d *Downloader) download(segIndex int) error {
   tsFilename := tsFilename(segIndex)
   tsUrl := d.tsURL(segIndex)
   b, e := Get(tsUrl)
   if e != nil {
   return fmt.Errorf("request %s, %s", tsUrl, e.Error())
   }
   //noinspection GoUnhandledErrorResult
   defer b.Close()
   fPath := filepath.Join(d.tsFolder, tsFilename)
   fTemp := fPath + tsTempFileSuffix
   f, err := os.Create(fTemp)
   if err != nil {
   return fmt.Errorf("create file: %s, %s", tsFilename, err.Error())
   }
   bytes, err := ioutil.ReadAll(b)
   if err != nil {
   return fmt.Errorf("read bytes: %s, %s", tsUrl, err.Error())
   }
   sf := d.result.M3u8.Segments[segIndex]
   if sf == nil {
   return fmt.Errorf("invalid segment index: %d", segIndex)
   }
   key, ok := d.result.Keys[sf.KeyIndex]
   if ok && key != "" {
   bytes, err = AES128Decrypt(bytes, []byte(key),
   []byte(d.result.M3u8.Keys[sf.KeyIndex].IV))
   if err != nil {
   return fmt.Errorf("decryt: %s, %s", tsUrl, err.Error())
   }
   }
   // https://en.wikipedia.org/wiki/MPEG_transport_stream
   // Some TS files do not start with SyncByte 0x47, they can not be played
   // after merging, Need to remove the bytes before the SyncByte 0x47(71).
   syncByte := uint8(71) //0x47
   bLen := len(bytes)
   for j := 0; j < bLen; j++ {
   if bytes[j] == syncByte {
   bytes = bytes[j:]
   break
   }
   }
   w := bufio.NewWriter(f)
   if _, err := w.Write(bytes); err != nil {
   return fmt.Errorf("write to %s: %s", fTemp, err.Error())
   }
   // Release file resource to rename file
   _ = f.Close()
   if err = os.Rename(fTemp, fPath); err != nil {
   return err
   }
   // Maybe it will be safer in this way...
   atomic.AddInt32(&d.finish, 1)
   fmt.Printf("[download %6.2f%%] %s\n", float32(d.finish)/float32(d.segLen)*100, tsUrl)
   return nil
}

////////////////////////////////////////////////////////////////////////////////

const (
   CryptMethodAES  CryptMethod = "AES-128"
   CryptMethodNONE CryptMethod = "NONE"
)

// #EXT-X-KEY:METHOD=AES-128,URI="key.key"
type Key struct {
	// 'AES-128' or 'NONE'
	// If the encryption method is NONE, the URI and the IV attributes MUST NOT be present
	Method CryptMethod
	URI    string
	IV     string
}


func FromURL(link string) (*Result, error) {
	u, err := url.Parse(link)
	if err != nil {
		return nil, err
	}
	link = u.String()
	body, err := Get(link)
	if err != nil {
		return nil, fmt.Errorf("request m3u8 URL failed: %s", err.Error())
	}
	//noinspection GoUnhandledErrorResult
	defer body.Close()
	m3u8, err := parse(body)
	if err != nil {
		return nil, err
	}
	if len(m3u8.MasterPlaylist) != 0 {
		sf := m3u8.MasterPlaylist[0]
		return FromURL(ResolveURL(u, sf.URI))
	}
	if len(m3u8.Segments) == 0 {
		return nil, errors.New("can not found any TS file description")
	}
	result := &Result{
		URL:  u,
		M3u8: m3u8,
		Keys: make(map[int]string),
	}

	for idx, key := range m3u8.Keys {
		switch {
		case key.Method == "" || key.Method == CryptMethodNONE:
			continue
		case key.Method == CryptMethodAES:
			// Request URL to extract decryption key
			keyURL := key.URI
			keyURL = ResolveURL(u, keyURL)
			resp, err := Get(keyURL)
			if err != nil {
				return nil, fmt.Errorf("extract key failed: %s", err.Error())
			}
			keyByte, err := ioutil.ReadAll(resp)
			_ = resp.Close()
			if err != nil {
				return nil, err
			}
			fmt.Println("decryption key: ", string(keyByte))
			result.Keys[idx] = string(keyByte)
		default:
			return nil, fmt.Errorf("unknown or unsupported cryption method: %s", key.Method)
		}
	}
	return result, nil
}

func parse(reader io.Reader) (*M3u8, error) {
   s := bufio.NewScanner(reader)
   var lines []string
   for s.Scan() {
      lines = append(lines, s.Text())
   }
   var (
      i     = 0
      count = len(lines)
      m3u8  = &M3u8{
         Keys: make(map[int]*Key),
      }
      keyIndex = 0
      key     *Key
      seg     *Segment
      extInf  bool
      extByte bool
   )
   for ; i < count; i++ {
      line := strings.TrimSpace(lines[i])
      if i == 0 {
         if line != "#EXTM3U" {
            return nil, fmt.Errorf("invalid m3u8, missing #EXTM3U in line 1")
         }
         continue
      }
      switch {
      case line == "":
         continue
      case strings.HasPrefix(line, "#EXT-X-PLAYLIST-TYPE:"):
         if _, err := fmt.Sscanf(line, "#EXT-X-PLAYLIST-TYPE:%s", &m3u8.PlaylistType); err != nil {
            return nil, err
         }
         isValid := m3u8.PlaylistType == "" || m3u8.PlaylistType == PlaylistTypeVOD || m3u8.PlaylistType == PlaylistTypeEvent
         if !isValid {
            return nil, fmt.Errorf("invalid playlist type: %s, line: %d", m3u8.PlaylistType, i+1)
         }
      case strings.HasPrefix(line, "#EXT-X-TARGETDURATION:"):
         if _, err := fmt.Sscanf(line, "#EXT-X-TARGETDURATION:%f", &m3u8.TargetDuration); err != nil {
            return nil, err
         }
      case strings.HasPrefix(line, "#EXT-X-MEDIA-SEQUENCE:"):
         if _, err := fmt.Sscanf(line, "#EXT-X-MEDIA-SEQUENCE:%d", &m3u8.MediaSequence); err != nil {
            return nil, err
         }
      case strings.HasPrefix(line, "#EXT-X-VERSION:"):
         if _, err := fmt.Sscanf(line, "#EXT-X-VERSION:%d", &m3u8.Version); err != nil {
            return nil, err
         }
      // Parse master playlist
      case strings.HasPrefix(line, "#EXT-X-STREAM-INF:"):
         mp, err := parseMasterPlaylist(line)
         if err != nil {
            return nil, err
         }
         i++
         mp.URI = lines[i]
         if mp.URI == "" || strings.HasPrefix(mp.URI, "#") {
            return nil, fmt.Errorf("invalid EXT-X-STREAM-INF URI, line: %d", i+1)
         }
         m3u8.MasterPlaylist = append(m3u8.MasterPlaylist, mp)
         continue
      case strings.HasPrefix(line, "#EXTINF:"):
         if extInf {
            return nil, fmt.Errorf("duplicate EXTINF: %s, line: %d", line, i+1)
         }
         if seg == nil {
            seg = new(Segment)
         }
         var s string
         if _, err := fmt.Sscanf(line, "#EXTINF:%s", &s); err != nil {
            return nil, err
         }
         if strings.Contains(s, ",") {
            split := strings.Split(s, ",")
            seg.Title = split[1]
            s = split[0]
         }
         df, err := strconv.ParseFloat(s, 32)
         if err != nil {
            return nil, err
         }
         seg.Duration = float32(df)
         seg.KeyIndex = keyIndex
         extInf = true
      case strings.HasPrefix(line, "#EXT-X-BYTERANGE:"):
         if extByte {
            return nil, fmt.Errorf("duplicate EXT-X-BYTERANGE: %s, line: %d", line, i+1)
         }
         if seg == nil {
            seg = new(Segment)
         }
         var b string
         if _, err := fmt.Sscanf(line, "#EXT-X-BYTERANGE:%s", &b); err != nil {
            return nil, err
         }
         if b == "" {
            return nil, fmt.Errorf("invalid EXT-X-BYTERANGE, line: %d", i+1)
         }
         if strings.Contains(b, "@") {
            split := strings.Split(b, "@")
            offset, err := strconv.ParseUint(split[1], 10, 64)
            if err != nil {
               return nil, err
            }
            seg.Offset = uint64(offset)
            b = split[0]
         }
         length, err := strconv.ParseUint(b, 10, 64)
         if err != nil {
            return nil, err
         }
         seg.Length = uint64(length)
         extByte = true
      // Parse segments URI
      case !strings.HasPrefix(line, "#"):
         if extInf {
            if seg == nil {
               return nil, fmt.Errorf("invalid line: %s", line)
            }
            seg.URI = line
            extByte = false
            extInf = false
            m3u8.Segments = append(m3u8.Segments, seg)
            seg = nil
            continue
         }
      // Parse key
      case strings.HasPrefix(line, "#EXT-X-KEY"):
         params := parseLineParameters(line)
         if len(params) == 0 {
            return nil, fmt.Errorf("invalid EXT-X-KEY: %s, line: %d", line, i+1)
         }
         method := CryptMethod(params["METHOD"])
         if method != "" && method != CryptMethodAES && method != CryptMethodNONE {
            return nil, fmt.Errorf("invalid EXT-X-KEY method: %s, line: %d", method, i+1)
         }
         keyIndex++
         key = new(Key)
         key.Method = method
         key.URI = params["URI"]
         key.IV = params["IV"]
         m3u8.Keys[keyIndex] = key
      case line == "#EndList":
         m3u8.EndList = true
      default:
         continue
      }
   }
   return m3u8, nil
}

