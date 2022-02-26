# ProtoBuf

I reject the idea of having to use a compiler for ProtoBuf. I think you should
be able to Marshal and Unmarshal just like JSON. And really, that should be
possible. If they had only added a single extra wiretype for messages, ProtoBuf
would be more or less self describing. At any rate, I want a package that can
decode arbitrary ProtoBuf, and can also encode some Map or Struct into ProtoBuf
as well.

- https://github.com/golang/protobuf/issues/1370
- https://stackoverflow.com/questions/41348512/protobuf-unmarshal-unknown

> Is it a string, bytes, or a sub-type? you don't know. You might be able to
> figure it out for a specific input, like when reverse engineering a gRPC api,
> but not in general.

You can generalize this. `string`, `bytes` and `message` all get passed as wire
type 2 (length-delimited). But you **can** differ between strings and bytes
[1], and not all bytes slices are valid messages. So you can run those tests,
and if you still have overlap, then you can just parse the data as both types,
and add both types to the output under the same field number, using the type as
the discriminator. This wouldnt work with `protowire.Type` [2], as again it
uses the same type for all three, so any implementation would need to create a
new `string`, `bytes` and `message` type.

1. https://github.com/golang/go/blob/go1.17.6/src/net/http/sniff.go#L297-L309
2. https://godocs.io/google.golang.org/protobuf/encoding/protowire#Type

> what's the wire type?

https://developers.google.com/protocol-buffers/docs/encoding#structure

> Why would you even want it vs using something else like JSON or BSON?

If it was my choice, I would never use ProtoBuf ever again. Its an awful
format. However some servers I deal with, require ProtoBuf request body, and
return ProtoBuf response body.

> why don't you use this plus some `protowire.EncodeTag`

Thats a good idea, but in my case I wanted an implementation that treats the
field name as first class citizen, so I ended up doing something like this:

~~~go
type Tag struct {
   protowire.Number
   Name string
}
~~~

Then I can use the `string` either as a type discriminator, or as the field
name associated with the field number.
