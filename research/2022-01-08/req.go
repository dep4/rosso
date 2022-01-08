package main

import (
   "github.com/89z/format/crypto"
   "io"
   "net/http"
   "net/http/httputil"
   "net/url"
   "os"
)

var req0 = &http.Request{Method:"POST", URL:&url.URL{Scheme:"https", Opaque:"",
User:(*url.Userinfo)(nil), Host:"music.amazon.com",
Path:"/NA/api/dmls/getDashManifestsV2", RawPath:"", ForceQuery:false,
RawQuery:"", Fragment:"", RawFragment:""},
Header:http.Header{"Accept":[]string{"application/json, text/javascript, */*"},
"Connection":[]string{"Keep-Alive"},
"Content-Length":[]string{"484"}, "Content-Type":[]string{"application/json; charset=UTF-8"}, "Host":[]string{"music.amazon.com"},
"User-Agent":[]string{"Harley/3.9.1.286 A1DL2DVDQVK3Q/17.19.4"},
"X-Adp-Alg":[]string{"SHA256withRSA:1.0"},
"X-Adp-Signature":[]string{"ONkEDeqaEbmmzCifn19jahB2FcKOgvDPdl0lQvHkW5UFfovf9cmiLele6FUdX2/Rw7zGEYad+dszadqX1aZ5X+poMdTi+BdhMlZUJOrdf7MwPC6ngCAh4fG1KeNAcBds8zkA7/kJNJQsn7g11Z7wXBAtkcUe6M+wksncPgzqMZWimmVBDGRHACRkSZM5+AG3b8V8/jvFKT+tDt1Etwb5JNpkpB2eYKELar9voB0EAe4tflrjZ+RQDrsgDgaCWhTnk8Mqtrr0vUvlLpd/vriWUA9BBOB7F5yBjslSThJWARvqVcnlDH6t3rQcR20GxKT0ZKivvwea2NBcpNTEZFr76w==:2022-01-08T15:30:30Z"},
"X-Adp-Token":[]string{"{enc:G8CuVtRvTPnv+OpI+JxO/formyxUgdR+vvtEe7kTY1dZbnu2BGgE2RS7zJn5w4pxPjPil+Bc557r3ppTUlYMOnwVlo8/ht0/45zO4EdbbGiAaN637eanjQqDoeV3nY6yc+1xAXBHgCfcDDZrbnrE+vwyHFa26T/5vCpjFE0WNKSizb5QWucOua7fObNq6p6HigkKx6ECNqw+FILeBtkZH9g3Wk6bb/MxM4fFvwlqu400rfwXPVTba/Z+pWRupf4Wuz14hCxjjfO4HxG41wbXG1y+ntKGH74jToMXgPmE0JaXMyqTdb9PPzveyhS3rPgnFxi0LsCCutqOcQOY72RN8l04iE28fN56DWewk2Yy30z5XtJKzyfGktpVEVUzNzQkGe+kryShCv9akoYbRuSyIQ6kMJK1F0YBgiyEuGB82ztWu8LeWEyLCd+OaNUzyhU3UkdsNFag2NHwskKm0vpH1uKxCh+yh+Fj9QFGEEd5s5FkYhAO8CA7qxE5lPsvGzKS/fHkZnUOjuS6Dv6XaIX5CrcaTrMsEvy4YTinS8iiRRZ3E1Wj27FVvyAxFnZFblns8KO8pI5hQn95vpDZYjaUOfU56waes0+c6M1rM/H5gWOdu2RNWm/4zvDlwPlVwP93n+xbW5E654FcuaGkgpzoP3/fUPAYvlGJip31tEqQxm9AFzM2zcD7JSsrWhXtOriRN6z4TL9y4wZFcwpwOy51OeZLg6EZqfL/QpzMzmC1SZdECGR7PNszUs+zZ45f2QnvKaqF8pjZpOGTxj/YFG/wiXTvSd2Ae29ghLW/Sei+OSgHPiIDR32nTZdspVooNvEjrvbqiOla00h+z8PeY7P1L2jJm+wm1SriDB89f155sHzHen6ycIPl5BH2W6mOsmH1yc7qSZN39F/54Dp3l36z7tK1FGeZIH36py4NUfQ7eBCG6G8yjlXfBeCjbDeRNn2My2Bic3GGrM0Wgo/0iaevJMQLQCh2FpImmBSnomcbkYIwQOy/YSZ/Bb4ufptwOGdMOCBaiGX9HOgfblI136+9fYymVVyQo6f/4TdZrER8lzA=}{key:ouP5kcb/rsDZSZAHHsZLYQsLBZE7ULJZiRpqe2rCTdSC+lyiH/Aq2VshAyr1SBPIbUZZoNJBtFJieDlayrerliOC0U5tkKsQ0gbPNHzpuWWKjrk7KQwUGCgoK0ZWDDBQJIyEnKDpmLBeDgKcX3I/cFE0zRU/q9UGNixZ/caBeLxh8CbVh5P4VJeHS2oLjXIIu/nbLjsr4B0QgOb1tuSSD7pQ0vlA2RLho1zAK0dJ6tNhRX8/sdXXeP7fyxKm0DfRYy6eeVuYTS4VVZaTTQcty1SSFlNwZ2E6fhfeJJT2oQnoTn+td6Ocmcdi975EvIi/J8mvBtk4xLmzblw6TckGOQ==}{iv:G15QPbhxm1DkGyuFlu/OyA==}{name:QURQVG9rZW5FbmNyeXB0aW9uS2V5}{serial:Mg==}"},
"X-Amz-Requestid":[]string{"d6661bfb-6cb8-4c30-bf61-22b5d68f033d"},
"X-Amz-Target":[]string{"com.amazon.digitalmusiclocator.DigitalMusicLocatorServiceExternal.getDashManifestsV2"}},
Body:io.NopCloser(body0)}

func main() {
   hello, err := crypto.ParseJA3(crypto.AndroidAPI26)
   if err != nil {
      panic(err)
   }
   res, err := crypto.Transport(hello).RoundTrip(req0)
   if err != nil {
      panic(err)
   }
   defer res.Body.Close()
   buf, err := httputil.DumpResponse(res, true)
   if err != nil {
      panic(err)
   }
   os.Stdout.Write(buf)
}
