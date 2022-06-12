package net

import (
   "net/url"
)

type Values url.Values

func (v Values) Get(key string) string {
   return url.Values(v).Get(key)
}

func (v *Values) UnmarshalText(text []byte) error {
   query, err := url.ParseQuery(string(text))
   if err != nil {
      return err
   }
   *v = Values(query)
   return nil
}
