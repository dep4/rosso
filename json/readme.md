# JSON

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
