Is it possible to generalize a function like this:

```go
type String string
func (String) name() string { return "String" }
type strings []String

func (s strings) get() *String {
   if len(s) == 0 { return nil }
   return &s[0]
}
```

here is what I tried: https://go.dev/play/p/wDnNIK9_uZb
