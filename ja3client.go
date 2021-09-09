package ja3

import (
   "github.com/refraction-networking/utls"
   "io"
   "net/http"
   "net/url"
   "strings"
)

// JA3Client contains is similar to http.Client
type JA3Client struct {
   *http.Client
   Config  *tls.Config
   Browser Browser
}

// New creates a JA3Client based on a Browser struct
func New(b Browser) (*JA3Client, error) {
   client, err := NewWithString(b.JA3)
   if err != nil {
      return nil, err
   }
   client.Browser = b
   return client, nil
}

// NewWithString creates a JA3 client with the specified JA3 string
func NewWithString(ja3 string) (*JA3Client, error) {
   tr, err := NewTransport(ja3)
   if err != nil {
      return nil, err
   }
   client := &http.Client{Transport: tr}
   return &JA3Client{
      client,
      &tls.Config{},
      Browser{JA3: ja3},
   }, nil
}

// Do sends an HTTP request and returns an HTTP response, following policy
// (such as redirects, cookies, auth) as configured on the client.
func (c *JA3Client) Do(req *http.Request) (*http.Response, error) {
	if _, ok := req.Header["User-Agent"]; !ok && c.Browser.UserAgent != "" {
		req.Header.Set("User-Agent", c.Browser.UserAgent)
	}

	return c.Client.Do(req)
}

// Get issues a GET to the specified URL.
func (c *JA3Client) Get(targetURL string) (*http.Response, error) {
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Post issues a POST to the specified URL.
func (c *JA3Client) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return c.Do(req)
}

// Head issues a HEAD to the specified URL.
func (c *JA3Client) Head(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// PostForm issues a POST to the specified URL,
// with data's keys and values URL-encoded as the request body.
func (c *JA3Client) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	return c.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

// Browser represents a browser JA3 and User-Agent string
type Browser struct {
	JA3       string
	UserAgent string
}

// ChromeAuto mocks Chrome 78
var ChromeAuto = Browser{
	JA3:       "769,47–53–5–10–49161–49162–49171–49172–50–56–19–4,0–10–11,23–24–25,0",
	UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.97 Safari/537.36",
}

// SafariAuto mocks Safari 604.1
var SafariAuto = Browser{
	JA3:       "771,4865-4866-4867-49196-49195-49188-49187-49162-49161-52393-49200-49199-49192-49191-49172-49171-52392-157-156-61-60-53-47-49160-49170-10,65281-0-23-13-5-18-16-11-51-45-43-10-21,29-23-24-25,0",
	UserAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 13_1_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.1 Mobile/15E148 Safari/604.1",
}
