package curlize

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"unicode/utf8"

	shellquote "github.com/kballard/go-shellquote"
)

var ErrNonUTF8Body = errors.New("request body is not utf-8 string")

type Command []string

func Curlize(r *http.Request) (Command, error) {
	var command []string

	command = append(command, "curl", "-X", r.Method)

	if r.Body != nil {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		if err := r.Body.Close(); err != nil {
			return nil, err
		}
		r.Body = ioutil.NopCloser(bytes.NewReader(body))
		if !utf8.Valid(body) {
			return nil, ErrNonUTF8Body
		}
		if len(body) != 0 {
			command = append(command, "-d", string(body))
		}
	}

	var keys []string
	for k := range r.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		command = append(command, "-H", fmt.Sprintf("%s: %s", k, strings.Join(r.Header[k], " ")))
	}

	command = append(command, r.URL.String())

	return command, nil
}

func (c Command) String() string {
	return shellquote.Join(c...)
}
