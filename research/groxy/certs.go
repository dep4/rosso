package groxy

import (
   "crypto/tls"
   "net/http"
)

var groxyCa tls.Certificate

func init() {
	var err error
	groxyCa, err = tls.X509KeyPair(caCert, caKey)
	if err != nil {
		panic(err)
	}
}

var caKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEApzPRppQb/lr92PAyn63Tu1HIcsO1F5FjlMRWzSHTjyUUtKs9
U171CHRS/91s8RKeP9yAFh3SRIKsV6lNrv4lifgRJguTklIS3XaYH0lsQ4OXc4tV
BRDWewFWV/BPa1ym3FVN1Iu4OXgkshEDIA5eP5PWWwrr0hxDG17Vs1s80nc8UB8D
dy0ai3hrzQQVWiu/CqOB6czd8SGqQqtUHM2N9cHRWv4OIWRlw0tS2J+vf6Jd0xuD
y0mKwsqK3vpKPaWXqnuq9MT6M+FkjZx9tSLB3TrOviUn5674Oyp59zxIWdILDXaN
Tu0xj5VRa8CuNPYoW91KPF+BxZMb/+QyxYVGTwIDAQABAoIBAQCcUj44V1C0xa1/
HVK3J3VFNHkLkx3EIxHPDKF6t51rv2dUYqS7RZQhi1/uB77KxHVfj2/RPaBQnTsz
2f3fFY1TKLIft8MIkeNBWpdu6N5nYKhARov+aHeeGOn43Zvi7IS0iqxxgw+B62mj
cXRSjBxhpH5MMDG6BuJWvJtfTHXY7k57+wxxUgim8jrb5ycrwoA26RJOzO1ZAQKi
kmayLTI27Purqdp1vSxjCM2BUKtmMvY94SEPr6Krnnu+QEX5lBoB0uKcnCificvT
7Chu/lAvSAxuE5qzdOceYO4IhpDGos9FqewF527iWwnJPnMY+WCauQuBNs98bo2d
J7twRI1BAoGBANEE86PzRWmA9qKQfCZjyU4zwAP+YSB2QqpRvvLdguvh4cuJX9Q7
/t+K0qE1BeQ0LGmVXOMYzKTAugeEF2C7BhOOHcEWCl5VnrIqnvQUdmnKpk9uwsEo
m9XBCDk55UDtoj7mrCI07tRyN1fam+mqBFpFiywDH04rwzTdUEozSoELAoGBAMzI
tNBWwS2T1k9HBJhHqp8gc3sDH0fM1dwPuEyBOOZQT3MujJXzEZ6NjJRRFS2zr5fa
CpK45wuqBNm8+J7ddn0qyTA7L9dILGNybsMr/iNFhLSQQvqWf6/70lDgryGSixbS
BUWzMrhj4p7NCi1MWtGlgWNWJ0wWbX3ulHlmiyJNAoGAdN0R497GNmDWhLnH0CcG
jGS+vPzjDSVRzKx72IolAwF+HzCllaIdtJqHfX6J1redBUNvdcGN2aev2zftYjXv
Bcv1stlB3tB8NB1EVi+CrU+SgOLqnNi5mF+e23AVT6INjqGmKFH0Hm/lpYcimRhn
6pjrYSY1wJ0TPWFn3LSkuVUCgYBZMRGPrDl3IGy3Grdlm9E0fW3Opg15uD1tG2up
5p5hsZjpEd0pcjS6WexgrPAMd7aC7PSt8tquud5i92tRxiNVYM0/tIehvE2ZAr47
Q23s7tpdyndhrSrv+z4e+71LYggyaEbjlKkCpn6Nq8BC+n6T0QIJwfxbp2kI8GC6
f28aAQKBgCtwuMfPjZMhB8OSiGRG2g/p6KfFXrTHxb85EjJXFL4Ucnz1Xy5UD14m
jnbySrCJnndFY1j4QiSqhoQDDi92+xr5tB9YU3DcXtxMfWGu/KxwbpDdlM3Inmva
Pt6+mU/Y5YkvlROaWI+DCMGzKdTg8V54rU71zcyOGLhxMr0Gvyls
-----END RSA PRIVATE KEY-----`)

var caCert = []byte(`-----BEGIN CERTIFICATE-----
MIIDtTCCAp2gAwIBAgIJANju9bI0gZV7MA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV
BAYTAkFVMRMwEQYDVQQIEwpTb21lLVN0YXRlMSEwHwYDVQQKExhJbnRlcm5ldCBX
aWRnaXRzIFB0eSBMdGQwHhcNMTcwNDA1MTQ1NDMyWhcNMTcwNzA0MTQ1NDMyWjBF
MQswCQYDVQQGEwJBVTETMBEGA1UECBMKU29tZS1TdGF0ZTEhMB8GA1UEChMYSW50
ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIB
CgKCAQEApzPRppQb/lr92PAyn63Tu1HIcsO1F5FjlMRWzSHTjyUUtKs9U171CHRS
/91s8RKeP9yAFh3SRIKsV6lNrv4lifgRJguTklIS3XaYH0lsQ4OXc4tVBRDWewFW
V/BPa1ym3FVN1Iu4OXgkshEDIA5eP5PWWwrr0hxDG17Vs1s80nc8UB8Ddy0ai3hr
zQQVWiu/CqOB6czd8SGqQqtUHM2N9cHRWv4OIWRlw0tS2J+vf6Jd0xuDy0mKwsqK
3vpKPaWXqnuq9MT6M+FkjZx9tSLB3TrOviUn5674Oyp59zxIWdILDXaNTu0xj5VR
a8CuNPYoW91KPF+BxZMb/+QyxYVGTwIDAQABo4GnMIGkMB0GA1UdDgQWBBSliIeu
1dMMozWayyJHvV5N9vPF9zB1BgNVHSMEbjBsgBSliIeu1dMMozWayyJHvV5N9vPF
96FJpEcwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgTClNvbWUtU3RhdGUxITAfBgNV
BAoTGEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZIIJANju9bI0gZV7MAwGA1UdEwQF
MAMBAf8wDQYJKoZIhvcNAQELBQADggEBAIXKZt4pHJuWX66+4l8ZXkeZvP5FeUXt
/b6dBgFVaId8m909HL2/eXUDyz6WPfCz+HhKQYCvqAZhwtYDFg24p++gFjWME2dh
sHPuZoee/RE0An7gN8lg75Ga2s4XBcWPabT03+gEwOXwr1GHvEto/+S5rxbQHgP6
FQ+OtuIpSAfuuJjkYdOzgaqxWKBu/tYJI+pgt5Xly+46Q0a2ovUTg4ff5GrA0VfT
bJNefkZ2L79jEO6aR0t/+hWgaM4XG++cgt6COU/ljzgFNOe8U7GJ0mL5keX2VuFP
WvcLOt83/KZ1jHrn5wkv0ajqtbJYXHu+e2kD3yoElZGxKTVJRSfV6/0=
-----END CERTIFICATE-----`)


// Handler handles http.Request and somehow generate http.Response or error.
type Handler func(*http.Request) (*http.Response, error)

// Middleware wraps original Handler and create new Handler.
type Middleware func(Handler) Handler

var httpclient = &http.Client{
	CheckRedirect: func(r *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

// DefaultHTTPHandler pass the request to the target server, and returns its response or error.
func DefaultHTTPHandler(req *http.Request) (*http.Response, error) {
	return httpclient.Do(req)
}

// DefaultHTTPSHandler pass the request to the target server, and returns its response or error.
func DefaultHTTPSHandler(tr *http.Transport) Handler {
	return tr.RoundTrip
}

// Logger is an interface to log events.
type Logger interface {
	Print(...interface{})
}

type nullLogger struct{}

func (nullLogger) Print(...interface{}) {}

// FuncLogger is a Logger that wraps print function
type FuncLogger func(...interface{})

// Print invokes f with args
func (f FuncLogger) Print(args ...interface{}) {
	f(args...)
}
