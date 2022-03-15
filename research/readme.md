# Research

- https://github.com/89z/format/blob/9636151e/protobuf/public.go#L11-L15
- https://godocs.io/google.golang.org/protobuf/testing/protopack#Message

If we have one value, that looks like this:

~~~go
message{
   token{1, 2, "hello"},
}
~~~

what if we have two values? It could be this:

~~~go
message{
   token{1, 2, "hello"},
   token{1, 2, "world"},
}
~~~

or this:

~~~go
message{
   token{
      1, 2, []string{"hello", "world"},
   },
}
~~~
