# Protobuf

- https://github.com/golang/protobuf/issues/1370
- https://stackoverflow.com/questions/26744873/converting-map-to-struct
- https://stackoverflow.com/questions/41348512/protobuf-unmarshal-unknown

I think I found a fix for ambiguous data. From my testing, problems only happen
with `protowire.BytesType`, as the result can be a `string`, `bytes` or
`message`. You can solve some cases by seeing if the data will parse as a
message, or by checking if the data is binary [1]. However if the data can
parse as a message, the result could be multiple types.

To solve this, I wrote a package that has separate types for message, string
and bytes. If I encounter ambiguous data, I add two entries. One as a string
(or bytes) and one as a message. I give both the same field number, and use my
new types as the discriminator.

1. https://github.com/golang/go/blob/go1.17.6/src/net/http/sniff.go#L297-L309
