# Protobuf

- bytes to map
- map to bytes
- map to map
- map to struct
- struct to map

We can do this:

~~~
bytes
map[protowire.Number]interface{}
struct
~~~

We need to do this:

~~~
struct
map[string]interface{}
map[protowire.Number]interface{}
bytes
~~~

- https://github.com/golang/protobuf/issues/1370
- https://github.com/philpearl/plenc/issues/3
- https://github.com/protocolbuffers/protobuf-go/blob/master/testing/protopack/pack.go
- https://github.com/segmentio/encoding/issues/103
- https://stackoverflow.com/questions/26744873/converting-map-to-struct
- https://stackoverflow.com/questions/41348512/protobuf-unmarshal-unknown
