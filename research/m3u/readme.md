# Research

Currently the Master type looks like this:

~~~go
type Master struct {
   Resolution string
   Bandwidth int64
   Codecs string
   URI string
}
~~~

If the URI is relative, then consumers of this API will need to prepend to the
path, which is fine. Currently the Segment type looks like this:

~~~go
type Segment struct {
   Key string
   URI []string
}
~~~

If the URI is relative, then consumers of this API will need to prepend to the
path. Is that OK?

https://stackoverflow.com/questions/27918208/go-get-parent-struct
