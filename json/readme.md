# JSON

14. https://github.com/goldeneggg/structil/issues/35
13. https://github.com/WillAbides/rjson/issues/47
12. https://github.com/segmentio/encoding/issues/99
11. https://github.com/tamerh/jsparser/issues/7

---

13. https://github.com/tidwall/pjson/issues/1
12. https://github.com/clarketm/json/issues/4
11. https://github.com/valyala/fastjson/issues/73
10. https://github.com/minio/simdjson-go/issues/43
9. https://github.com/ohler55/ojg/issues/76
8. https://github.com/Jeffail/gabs/issues/110
7. https://github.com/json-iterator/go/issues/577
6. https://github.com/goccy/go-json/issues/292
5. https://github.com/pquerna/ffjson/issues/263
4. https://github.com/pkg/json/issues/12
3. https://github.com/mailru/easyjson/issues/351
2. https://github.com/tidwall/gjson/issues/235
1. https://github.com/buger/jsonparser/issues/235

~~~
json language:go pushed:>2020-09-28 stars:>5
~~~

## Non strict mode

I am trying to parse JSON from some input like this:

~~~js
window.__additionalDataLoaded('extra',{"shortcode_media":999});
~~~

Where the output would be:

~~~json
{"shortcode_media":999}
~~~

is something like that possible? Similar to `Strict = false` with XML:

https://godocs.io/encoding/xml#Decoder.Strict
