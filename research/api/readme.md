# API

~~~
func (x *Index) Read(r io.Reader) error
func (x *Index) Write(w io.Writer) error
~~~

https://godocs.io/index/suffixarray

~~~
func ReadRequest(b *bufio.Reader) (*Request, error)
func (r *Request) Write(w io.Writer) error

func ReadResponse(r *bufio.Reader, req *Request) (*Response, error)
func (r *Response) Write(w io.Writer) error
~~~

https://godocs.io/net/http
