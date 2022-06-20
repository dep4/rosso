# API

## io.Reader

~~~
10:pkg crypto/ed25519, func GenerateKey(io.Reader) (PublicKey, PrivateKey, error)
120:pkg os, method (*File) ReadFrom(io.Reader) (int64, error)
12:pkg archive/tar, func NewReader(io.Reader) *Reader
12:pkg crypto, type Decrypter interface, Decrypt(io.Reader, []uint8, DecrypterOpts) ([]uint8, error)
1592:pkg debug/elf, func NewFile(io.ReaderAt) (*File, error)
163:pkg debug/plan9obj, func NewFile(io.ReaderAt) (*File, error)
16:pkg crypto/ed25519, method (PrivateKey) Sign(io.Reader, []uint8, crypto.SignerOpts) ([]uint8, error)
1726:pkg debug/elf, type Prog struct, embedded io.ReaderAt
173:pkg crypto/x509, func EncryptPEMBlock(io.Reader, string, []uint8, []uint8, PEMCipher) (*pem.Block, error)
1778:pkg debug/elf, type Section struct, embedded io.ReaderAt
17:pkg crypto/rsa, method (*PrivateKey) Decrypt(io.Reader, []uint8, crypto.DecrypterOpts) ([]uint8, error)
181:pkg debug/plan9obj, type Section struct, embedded io.ReaderAt
185:pkg bytes, method (*Buffer) ReadFrom(io.Reader) (int64, error)
1908:pkg debug/macho, func NewFile(io.ReaderAt) (*File, error)
1:pkg archive/zip, method (*File) OpenRaw() (io.Reader, error)
2044:pkg debug/macho, type Section struct, embedded io.ReaderAt
204:pkg net/http, func NewRequestWithContext(context.Context, string, string, io.Reader) (*Request, error)
2083:pkg debug/macho, type Segment struct, embedded io.ReaderAt
209:pkg compress/bzip2, func NewReader(io.Reader) io.Reader
211:pkg net/http/httptest, func NewRequest(string, string, io.Reader) *http.Request
2163:pkg debug/pe, func NewFile(io.ReaderAt) (*File, error)
216:pkg compress/flate, func NewReader(io.Reader) io.ReadCloser
217:pkg compress/flate, func NewReaderDict(io.Reader, []uint8) io.ReadCloser
2194:pkg debug/pe, type Section struct, embedded io.ReaderAt
21:pkg crypto/x509, func CreateCertificateRequest(io.Reader, *CertificateRequest, interface{}) ([]uint8, error)
2220:pkg encoding/ascii85, func NewDecoder(io.Reader) io.Reader
2249:pkg encoding/base32, func NewDecoder(*Encoding, io.Reader) io.Reader
2263:pkg encoding/base64, func NewDecoder(*Encoding, io.Reader) io.Reader
2282:pkg encoding/binary, func Read(io.Reader, ByteOrder, interface{}) error
2299:pkg encoding/csv, func NewReader(io.Reader) *Reader
22:pkg crypto/x509, func CreateRevocationList(io.Reader, *RevocationList, *Certificate, crypto.Signer) ([]uint8, error)
2325:pkg encoding/gob, func NewDecoder(io.Reader) *Decoder
2358:pkg encoding/json, func NewDecoder(io.Reader) *Decoder
2413:pkg encoding/xml, func NewDecoder(io.Reader) *Decoder
2438:pkg encoding/xml, type Decoder struct, CharsetReader func(string, io.Reader) (io.Reader, error)
243:pkg compress/gzip, func NewReader(io.Reader) (*Reader, error)
24:pkg bufio, func NewScanner(io.Reader) *Scanner
2582:pkg fmt, func Fscan(io.Reader, ...interface{}) (int, error)
2583:pkg fmt, func Fscanf(io.Reader, string, ...interface{}) (int, error)
2584:pkg fmt, func Fscanln(io.Reader, ...interface{}) (int, error)
264:pkg compress/lzw, func NewReader(io.Reader, Order, int) io.ReadCloser
26:pkg crypto/rsa, func SignPSS(io.Reader, *PrivateKey, crypto.Hash, []uint8, *PSSOptions) ([]uint8, error)
271:pkg compress/zlib, func NewReader(io.Reader) (io.ReadCloser, error)
272:pkg compress/zlib, func NewReaderDict(io.Reader, []uint8) (io.ReadCloser, error)
27:pkg crypto, type Signer interface, Sign(io.Reader, []uint8, SignerOpts) ([]uint8, error)
30251:pkg testing/iotest, func DataErrReader(io.Reader) io.Reader
30252:pkg testing/iotest, func HalfReader(io.Reader) io.Reader
30253:pkg testing/iotest, func NewReadLogger(string, io.Reader) io.Reader
30255:pkg testing/iotest, func OneByteReader(io.Reader) io.Reader
30256:pkg testing/iotest, func TimeoutReader(io.Reader) io.Reader
30300:pkg text/scanner, method (*Scanner) Init(io.Reader) *Scanner
31:pkg crypto/ecdsa, method (*PrivateKey) Sign(io.Reader, []uint8, crypto.SignerOpts) ([]uint8, error)
3410:pkg image, func Decode(io.Reader) (Image, string, error)
3411:pkg image, func DecodeConfig(io.Reader) (Config, string, error)
3425:pkg image, func RegisterFormat(string, string, func(io.Reader) (Image, error), func(io.Reader) (Config, error))
34:pkg crypto/rsa, method (*PrivateKey) Sign(io.Reader, []uint8, crypto.SignerOpts) ([]uint8, error)
35:pkg bufio, method (*Writer) ReadFrom(io.Reader) (int64, error)
360:pkg crypto/cipher, type StreamReader struct, R io.Reader
3683:pkg image/gif, func Decode(io.Reader) (image.Image, error)
3684:pkg image/gif, func DecodeAll(io.Reader) (*GIF, error)
3685:pkg image/gif, func DecodeConfig(io.Reader) (image.Config, error)
3691:pkg image/jpeg, func Decode(io.Reader) (image.Image, error)
3692:pkg image/jpeg, func DecodeConfig(io.Reader) (image.Config, error)
36:pkg bufio, method (ReadWriter) ReadFrom(io.Reader) (int64, error)
3703:pkg image/png, func Decode(io.Reader) (image.Image, error)
3704:pkg image/png, func DecodeConfig(io.Reader) (image.Config, error)
3714:pkg index/suffixarray, method (*Index) Read(io.Reader) error
3722:pkg io, func NewSectionReader(ReaderAt, int64, int64) *SectionReader
3735:pkg io, method (*SectionReader) Read([]uint8) (int, error)
3736:pkg io, method (*SectionReader) ReadAt([]uint8, int64) (int, error)
3737:pkg io, method (*SectionReader) Seek(int64, int) (int64, error)
3738:pkg io, method (*SectionReader) Size() int64
375:pkg crypto/dsa, func GenerateKey(*PrivateKey, io.Reader) error
376:pkg crypto/dsa, func GenerateParameters(*Parameters, io.Reader, ParameterSizes) error
3779:pkg io, type SectionReader struct
377:pkg crypto/dsa, func Sign(io.Reader, *PrivateKey, []uint8) (*big.Int, *big.Int, error)
3799:pkg io/ioutil, func NopCloser(io.Reader) io.ReadCloser
3800:pkg io/ioutil, func ReadAll(io.Reader) ([]uint8, error)
391:pkg crypto/ecdsa, func GenerateKey(elliptic.Curve, io.Reader) (*PrivateKey, error)
392:pkg crypto/ecdsa, func Sign(io.Reader, *PrivateKey, []uint8) (*big.Int, *big.Int, error)
3:pkg crypto/ecdsa, func SignASN1(io.Reader, *PrivateKey, []uint8) ([]uint8, error)
413:pkg crypto/elliptic, func GenerateKey(Curve, io.Reader) ([]uint8, *big.Int, *big.Int, error)
4323:pkg mime/multipart, func NewReader(io.Reader, string) *Reader
4447:pkg net, method (*TCPConn) ReadFrom(io.Reader) (int64, error)
444:pkg crypto/rand, func Int(io.Reader, *big.Int) (*big.Int, error)
445:pkg crypto/rand, func Prime(io.Reader, int) (*big.Int, error)
447:pkg crypto/rand, var Reader io.Reader
44:pkg archive/zip, func NewReader(io.ReaderAt, int64) (*Reader, error)
454:pkg crypto/rsa, func DecryptOAEP(hash.Hash, io.Reader, *PrivateKey, []uint8, []uint8) ([]uint8, error)
455:pkg crypto/rsa, func DecryptPKCS1v15(io.Reader, *PrivateKey, []uint8) ([]uint8, error)
456:pkg crypto/rsa, func DecryptPKCS1v15SessionKey(io.Reader, *PrivateKey, []uint8, []uint8) error
457:pkg crypto/rsa, func EncryptOAEP(hash.Hash, io.Reader, *PublicKey, []uint8, []uint8) ([]uint8, error)
458:pkg crypto/rsa, func EncryptPKCS1v15(io.Reader, *PublicKey, []uint8) ([]uint8, error)
459:pkg crypto/rsa, func GenerateKey(io.Reader, int) (*PrivateKey, error)
460:pkg crypto/rsa, func GenerateMultiPrimeKey(io.Reader, int, int) (*PrivateKey, error)
461:pkg crypto/rsa, func SignPKCS1v15(io.Reader, *PrivateKey, crypto.Hash, []uint8) ([]uint8, error)
4682:pkg net/http, func NewRequest(string, string, io.Reader) (*Request, error)
4687:pkg net/http, func Post(string, string, io.Reader) (*Response, error)
4691:pkg net/http, func ReadRequest(*bufio.Reader) (*Request, error)
4692:pkg net/http, func ReadResponse(*bufio.Reader, *Request) (*Response, error)
4705:pkg net/http, method (*Client) Post(string, string, io.Reader) (*Response, error)
479:pkg testing/iotest, func ErrReader(error) io.Reader
480:pkg testing/iotest, func TestReader(io.Reader, []uint8) error
4894:pkg net/http/httputil, func NewChunkedReader(io.Reader) io.Reader
4896:pkg net/http/httputil, func NewClientConn(net.Conn, *bufio.Reader) *ClientConn
4897:pkg net/http/httputil, func NewProxyClientConn(net.Conn, *bufio.Reader) *ClientConn
4898:pkg net/http/httputil, func NewServerConn(net.Conn, *bufio.Reader) *ServerConn
4902:pkg net/http/httputil, method (*ClientConn) Hijack() (net.Conn, *bufio.Reader)
4908:pkg net/http/httputil, method (*ServerConn) Hijack() (net.Conn, *bufio.Reader)
4927:pkg net/mail, func ReadMessage(io.Reader) (*Message, error)
4937:pkg net/mail, type Message struct, Body io.Reader
4:pkg compress/gzip, method (*Reader) Reset(io.Reader) error
5026:pkg net/textproto, func NewReader(*bufio.Reader) *Reader
5030:pkg net/textproto, method (*Conn) DotReader() io.Reader
5053:pkg net/textproto, method (*Reader) DotReader() io.Reader
5081:pkg net/textproto, type Reader struct, R *bufio.Reader
5303:pkg os/exec, type Cmd struct, Stdin io.Reader
548:pkg crypto/tls, type Config struct, Rand io.Reader
567:pkg encoding/hex, func NewDecoder(io.Reader) io.Reader
593:pkg crypto/x509, func CreateCertificate(io.Reader, *Certificate, *Certificate, interface{}, interface{}) ([]uint8, error)
5:pkg archive/zip, type Decompressor func(io.Reader) io.ReadCloser
5:pkg debug/buildinfo, func Read(io.ReaderAt) (*debug.BuildInfo, error)
610:pkg crypto/x509, method (*Certificate) CreateCRL(io.Reader, interface{}, []pkix.RevokedCertificate, time.Time, time.Time) ([]uint8, error)
66:pkg debug/macho, func NewFatFile(io.ReaderAt) (*FatFile, error)
6:pkg bufio, method (*Reader) Reset(io.Reader)
6:pkg compress/flate, type Resetter interface, Reset(io.Reader, []uint8) error
6:pkg compress/lzw, method (*Reader) Reset(io.Reader, Order, int)
850:pkg encoding/json, method (*Decoder) Buffered() io.Reader
851:pkg mime, type WordDecoder struct, CharsetReader func(string, io.Reader) (io.Reader, error)
853:pkg mime/quotedprintable, func NewReader(io.Reader) *Reader
88:pkg bufio, func NewReader(io.Reader) *Reader
89:pkg bufio, func NewReaderSize(io.Reader, int) *Reader
8:pkg compress/zlib, type Resetter interface, Reset(io.Reader, []uint8) error
~~~

## io.Writer

~~~
108:pkg image/gif, func Encode(io.Writer, image.Image, *Options) error
109:pkg image/gif, func EncodeAll(io.Writer, *GIF) error
10:pkg compress/zlib, method (*Writer) Reset(io.Writer)
13:pkg archive/tar, func NewWriter(io.Writer) *Writer
197:pkg bytes, method (*Buffer) WriteTo(io.Writer) (int64, error)
210:pkg net/http/cgi, type Handler struct, Stderr io.Writer
211:pkg net, method (*Buffers) WriteTo(io.Writer) (int64, error)
218:pkg compress/flate, func NewWriter(io.Writer, int) (*Writer, error)
219:pkg compress/flate, func NewWriterDict(io.Writer, int, []uint8) (*Writer, error)
2221:pkg encoding/ascii85, func NewEncoder(io.Writer) io.WriteCloser
2250:pkg encoding/base32, func NewEncoder(*Encoding, io.Writer) io.WriteCloser
2264:pkg encoding/base64, func NewEncoder(*Encoding, io.Writer) io.WriteCloser
2288:pkg encoding/binary, func Write(io.Writer, ByteOrder, interface{}) error
22:pkg log, method (*Logger) Writer() io.Writer
2300:pkg encoding/csv, func NewWriter(io.Writer) *Writer
2326:pkg encoding/gob, func NewEncoder(io.Writer) *Encoder
2346:pkg encoding/hex, func Dumper(io.Writer) io.WriteCloser
2359:pkg encoding/json, func NewEncoder(io.Writer) *Encoder
2402:pkg encoding/pem, func Encode(io.Writer, *Block) error
2410:pkg encoding/xml, func Escape(io.Writer, []uint8)
2414:pkg encoding/xml, func NewEncoder(io.Writer) *Encoder
244:pkg compress/gzip, func NewWriter(io.Writer) *Writer
245:pkg compress/gzip, func NewWriterLevel(io.Writer, int) (*Writer, error)
2555:pkg flag, method (*FlagSet) SetOutput(io.Writer)
2579:pkg fmt, func Fprint(io.Writer, ...interface{}) (int, error)
2580:pkg fmt, func Fprintf(io.Writer, string, ...interface{}) (int, error)
2581:pkg fmt, func Fprintln(io.Writer, ...interface{}) (int, error)
2633:pkg go/ast, func Fprint(io.Writer, *token.FileSet, interface{}, FieldFilter) error
2636:pkg strings, method (*Reader) WriteTo(io.Writer) (int64, error)
265:pkg compress/lzw, func NewWriter(io.Writer, Order, int) io.WriteCloser
273:pkg compress/zlib, func NewWriter(io.Writer) *Writer
274:pkg compress/zlib, func NewWriterLevel(io.Writer, int) (*Writer, error)
275:pkg compress/zlib, func NewWriterLevelDict(io.Writer, int, []uint8) (*Writer, error)
29:pkg bufio, method (*Reader) WriteTo(io.Writer) (int64, error)
30254:pkg testing/iotest, func NewWriteLogger(string, io.Writer) io.Writer
30257:pkg testing/iotest, func TruncateWriter(io.Writer, int64) io.Writer
30327:pkg text/tabwriter, func NewWriter(io.Writer, int, int, int, uint8, uint) *Writer
30329:pkg text/tabwriter, method (*Writer) Init(io.Writer, int, int, int, uint8, uint) *Writer
30332:pkg text/template, func HTMLEscape(io.Writer, []uint8)
30335:pkg text/template, func JSEscape(io.Writer, []uint8)
30346:pkg text/template, method (*Template) Execute(io.Writer, interface{}) error
30347:pkg text/template, method (*Template) ExecuteTemplate(io.Writer, string, interface{}) error
3099:pkg go/doc, func ToHTML(io.Writer, string, map[string]string)
3100:pkg go/doc, func ToText(io.Writer, string, string, string, int)
3154:pkg go/printer, func Fprint(io.Writer, *token.FileSet, interface{}) error
3155:pkg go/printer, method (*Config) Fprint(io.Writer, *token.FileSet, interface{}) error
3164:pkg go/scanner, func PrintError(io.Writer, error)
3368:pkg html/template, func HTMLEscape(io.Writer, []uint8)
3371:pkg html/template, func JSEscape(io.Writer, []uint8)
3383:pkg html/template, method (*Template) Execute(io.Writer, interface{}) error
3384:pkg html/template, method (*Template) ExecuteTemplate(io.Writer, string, interface{}) error
356:pkg image/png, method (*Encoder) Encode(io.Writer, image.Image) error
365:pkg crypto/cipher, type StreamWriter struct, W io.Writer
3693:pkg image/jpeg, func Encode(io.Writer, image.Image, *Options) error
3705:pkg image/png, func Encode(io.Writer, image.Image) error
3715:pkg index/suffixarray, method (*Index) Write(io.Writer) error
37:pkg bufio, method (ReadWriter) WriteTo(io.Writer) (int64, error)
3806:pkg io/ioutil, var Discard io.Writer
3817:pkg log, func New(io.Writer, string, int) *Logger
3826:pkg log, func SetOutput(io.Writer)
3:pkg archive/zip, method (*Writer) CreateRaw(*FileHeader) (io.Writer, error)
4324:pkg mime/multipart, func NewWriter(io.Writer) *Writer
4335:pkg mime/multipart, method (*Writer) CreateFormField(string) (io.Writer, error)
4336:pkg mime/multipart, method (*Writer) CreateFormFile(string, string) (io.Writer, error)
4337:pkg mime/multipart, method (*Writer) CreatePart(textproto.MIMEHeader) (io.Writer, error)
45:pkg archive/zip, func NewWriter(io.Writer) *Writer
4721:pkg net/http, method (*Request) Write(io.Writer) error
4722:pkg net/http, method (*Request) WriteProxy(io.Writer) error
4726:pkg net/http, method (*Response) Write(io.Writer) error
4742:pkg net/http, method (Header) Write(io.Writer) error
4743:pkg net/http, method (Header) WriteSubset(io.Writer, map[string]bool) error
47:pkg bytes, method (*Reader) WriteTo(io.Writer) (int64, error)
4895:pkg net/http/httputil, func NewChunkedWriter(io.Writer) io.WriteCloser
48:pkg crypto/tls, type Config struct, KeyLogWriter io.Writer
4:pkg archive/zip, type Compressor func(io.Writer) (io.WriteCloser, error)
5027:pkg net/textproto, func NewWriter(*bufio.Writer) *Writer
5083:pkg net/textproto, type Writer struct, W *bufio.Writer
5302:pkg os/exec, type Cmd struct, Stderr io.Writer
5304:pkg os/exec, type Cmd struct, Stdout io.Writer
55:pkg log, func Writer() io.Writer
568:pkg encoding/hex, func NewEncoder(io.Writer) io.Writer
5749:pkg runtime/pprof, func StartCPUProfile(io.Writer) error
5751:pkg runtime/pprof, func WriteHeapProfile(io.Writer) error
5756:pkg runtime/pprof, method (*Profile) WriteTo(io.Writer, int) error
575:pkg flag, method (*FlagSet) Output() io.Writer
585:pkg go/types, method (*Scope) WriteTo(io.Writer, int, bool)
5875:pkg strings, method (*Replacer) WriteString(io.Writer, string) (int, error)
60:pkg archive/zip, method (*Writer) Create(string) (io.Writer, error)
61:pkg archive/zip, method (*Writer) CreateHeader(*FileHeader) (io.Writer, error)
767:pkg log, method (*Logger) SetOutput(io.Writer)
7:pkg bufio, method (*Writer) Reset(io.Writer)
854:pkg mime/quotedprintable, func NewWriter(io.Writer) *Writer
857:pkg encoding/xml, func EscapeText(io.Writer, []uint8) error
883:pkg runtime/trace, func Start(io.Writer) error
8:pkg compress/flate, method (*Writer) Reset(io.Writer)
8:pkg compress/lzw, method (*Writer) Reset(io.Writer, Order, int)
901:pkg go/format, func Node(io.Writer, *token.FileSet, interface{}) error
90:pkg bufio, func NewWriter(io.Writer) *Writer
91:pkg bufio, func NewWriterSize(io.Writer, int) *Writer
9:pkg compress/gzip, method (*Writer) Reset(io.Writer)
~~~
