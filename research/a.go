package hello

type String string

func (String) name() string { return "String" }

type strings []String

func (s strings) get() *String {
   if len(s) == 0 {
      return nil
   }
   return &s[0]
}
