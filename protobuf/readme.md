# Protobuf

- https://github.com/golang/protobuf/issues/1370
- https://github.com/philpearl/plenc/issues/3
- https://github.com/protocolbuffers/protobuf-go/blob/master/testing/protopack/pack.go
- https://github.com/segmentio/encoding/issues/103
- https://stackoverflow.com/questions/26744873/converting-map-to-struct
- https://stackoverflow.com/questions/41348512/protobuf-unmarshal-unknown

## protobuf bytes to struct

1. protobuf `[]byte`
2. `map[protowire.Number]interface{}` DONE
3. JSON `[]byte` DONE
4. `struct` DONE

## struct to protobuf bytes

1. `struct`
2. JSON `[]byte` DONE
3. `map[string]interface{}` DONE

Then from three can do this:

~~~
map[protowire.Number]interface{} DONE
~~~

or can we go straight to this:

~~~
protobuf []byte
~~~
