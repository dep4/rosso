package net

import (
   "net/url"
   "strconv"
   "strings"
)

func ParseURL(s, sep string) (*url.URL, error) {
   _, after, found := strings.Cut(s, sep)
   if !found {
      return nil, notFound{sep}
   }
   addr, err := url.Parse(after)
   if err != nil {
      return nil, err
   }
   addr.RawQuery = ""
   return addr, nil
}

type notFound struct {
   value string
}

func (n notFound) Error() string {
   return strconv.Quote(n.value) + " is not found"
}
