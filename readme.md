# Parse

Go parsers for web formats

https://godocs.io/github.com/89z/parse

## HTML package

Takes HTML input, and can iterate through elements by tag name or by attribute
name and value. Content from text nodes can be returned. Also, you can check if
an element has a certain attribute, and return an attribute value given an
attribute name. Finally, you can indent and write the HTML to some output.

## JavaScript package

Takes JavaScript input, and will return a `map`. Keys are the variable names,
and values are the variable values. The values are returned as `byte` slices,
to make it easy to `json.Unmarshal`.
