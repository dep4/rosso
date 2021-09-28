// Package buffer contains buffer and wrapper types for byte slices. It is useful for writing lexers or other high-performance byte slice handling.
// The `Reader` and `Writer` types implement the `io.Reader` and `io.Writer` respectively and provide a thinner and faster interface than `bytes.Buffer`.
// The `Lexer` type is useful for building lexers because it keeps track of the start and end position of a byte selection, and shifts the bytes whenever a valid token is found.
// The `StreamLexer` does the same, but keeps a buffer pool so that it reads a limited amount at a time, allowing to parse from streaming sources.
package buffer

import (
	"io"
	"io/ioutil"
)

var nullBuffer = []byte{0}

// Lexer is a buffered reader that allows peeking forward and shifting, taking an io.Reader.
// It keeps data in-memory until Free, taking a byte length, is called to move beyond the data.
type Lexer struct {
	buf   []byte
	pos   int // index in buf
	start int // index in buf
	err   error

	restore func()
}

// NewLexer returns a new Lexer for a given io.Reader, and uses ioutil.ReadAll to read it into a byte slice.
// If the io.Reader implements Bytes, that is used instead.
// It will append a NULL at the end of the buffer.
func NewLexer(r io.Reader) *Lexer {
	var b []byte
	if r != nil {
		if buffer, ok := r.(interface {
			Bytes() []byte
		}); ok {
			b = buffer.Bytes()
		} else {
			var err error
			b, err = ioutil.ReadAll(r)
			if err != nil {
				return &Lexer{
					buf: nullBuffer,
					err: err,
				}
			}
		}
	}
	return NewLexerBytes(b)
}

// NewLexerBytes returns a new Lexer for a given byte slice, and appends NULL at the end.
// To avoid reallocation, make sure the capacity has room for one more byte.
func NewLexerBytes(b []byte) *Lexer {
	z := &Lexer{
		buf: b,
	}

	n := len(b)
	if n == 0 {
		z.buf = nullBuffer
	} else {
		// Append NULL to buffer, but try to avoid reallocation
		if cap(b) > n {
			// Overwrite next byte but restore when done
			b = b[:n+1]
			c := b[n]
			b[n] = 0

			z.buf = b
			z.restore = func() {
				b[n] = c
			}
		} else {
			z.buf = append(b, 0)
		}
	}
	return z
}

// Restore restores the replaced byte past the end of the buffer by NULL.
func (z *Lexer) Restore() {
	if z.restore != nil {
		z.restore()
		z.restore = nil
	}
}

// Err returns the error returned from io.Reader or io.EOF when the end has been reached.
func (z *Lexer) Err() error {
	return z.PeekErr(0)
}

// PeekErr returns the error at position pos. When pos is zero, this is the same as calling Err().
func (z *Lexer) PeekErr(pos int) error {
	if z.err != nil {
		return z.err
	} else if z.pos+pos >= len(z.buf)-1 {
		return io.EOF
	}
	return nil
}

// Peek returns the ith byte relative to the end position.
// Peek returns 0 when an error has occurred, Err returns the error.
func (z *Lexer) Peek(pos int) byte {
	pos += z.pos
	return z.buf[pos]
}

// PeekRune returns the rune and rune length of the ith byte relative to the end position.
func (z *Lexer) PeekRune(pos int) (rune, int) {
	// from unicode/utf8
	c := z.Peek(pos)
	if c < 0xC0 || z.Peek(pos+1) == 0 {
		return rune(c), 1
	} else if c < 0xE0 || z.Peek(pos+2) == 0 {
		return rune(c&0x1F)<<6 | rune(z.Peek(pos+1)&0x3F), 2
	} else if c < 0xF0 || z.Peek(pos+3) == 0 {
		return rune(c&0x0F)<<12 | rune(z.Peek(pos+1)&0x3F)<<6 | rune(z.Peek(pos+2)&0x3F), 3
	}
	return rune(c&0x07)<<18 | rune(z.Peek(pos+1)&0x3F)<<12 | rune(z.Peek(pos+2)&0x3F)<<6 | rune(z.Peek(pos+3)&0x3F), 4
}

// Move advances the position.
func (z *Lexer) Move(n int) {
	z.pos += n
}

// Pos returns a mark to which can be rewinded.
func (z *Lexer) Pos() int {
	return z.pos - z.start
}

// Rewind rewinds the position to the given position.
func (z *Lexer) Rewind(pos int) {
	z.pos = z.start + pos
}

// Lexeme returns the bytes of the current selection.
func (z *Lexer) Lexeme() []byte {
	return z.buf[z.start:z.pos:z.pos]
}

// Skip collapses the position to the end of the selection.
func (z *Lexer) Skip() {
	z.start = z.pos
}

// Shift returns the bytes of the current selection and collapses the position to the end of the selection.
func (z *Lexer) Shift() []byte {
	b := z.buf[z.start:z.pos:z.pos]
	z.start = z.pos
	return b
}

// Offset returns the character position in the buffer.
func (z *Lexer) Offset() int {
	return z.pos
}

// Bytes returns the underlying buffer.
func (z *Lexer) Bytes() []byte {
	return z.buf[: len(z.buf)-1 : len(z.buf)-1]
}

// Reset resets position to the underlying buffer.
func (z *Lexer) Reset() {
	z.start = 0
	z.pos = 0
}


// defaultBufSize specifies the default initial length of internal buffers.
var defaultBufSize = 4096

// MinBuf specifies the default initial length of internal buffers.
// Solely here to support old versions of parse.
var MinBuf = defaultBufSize


// Reader implements an io.Reader over a byte slice.
type Reader struct {
	buf []byte
	pos int
}

// NewReader returns a new Reader for a given byte slice.
func NewReader(buf []byte) *Reader {
	return &Reader{
		buf: buf,
	}
}

// Read reads bytes into the given byte slice and returns the number of bytes read and an error if occurred.
func (r *Reader) Read(b []byte) (n int, err error) {
	if len(b) == 0 {
		return 0, nil
	}
	if r.pos >= len(r.buf) {
		return 0, io.EOF
	}
	n = copy(b, r.buf[r.pos:])
	r.pos += n
	return
}

// Bytes returns the underlying byte slice.
func (r *Reader) Bytes() []byte {
	return r.buf
}

// Reset resets the position of the read pointer to the beginning of the underlying byte slice.
func (r *Reader) Reset() {
	r.pos = 0
}

// Len returns the length of the buffer.
func (r *Reader) Len() int {
	return len(r.buf)
}

// Writer implements an io.Writer over a byte slice.
type Writer struct {
	buf []byte
}

// NewWriter returns a new Writer for a given byte slice.
func NewWriter(buf []byte) *Writer {
	return &Writer{
		buf: buf,
	}
}

// Write writes bytes from the given byte slice and returns the number of bytes written and an error if occurred. When err != nil, n == 0.
func (w *Writer) Write(b []byte) (int, error) {
	n := len(b)
	end := len(w.buf)
	if end+n > cap(w.buf) {
		buf := make([]byte, end, 2*cap(w.buf)+n)
		copy(buf, w.buf)
		w.buf = buf
	}
	w.buf = w.buf[:end+n]
	return copy(w.buf[end:], b), nil
}

// Len returns the length of the underlying byte slice.
func (w *Writer) Len() int {
	return len(w.buf)
}

// Bytes returns the underlying byte slice.
func (w *Writer) Bytes() []byte {
	return w.buf
}

// Reset empties and reuses the current buffer. Subsequent writes will overwrite the buffer, so any reference to the underlying slice is invalidated after this call.
func (w *Writer) Reset() {
	w.buf = w.buf[:0]
}
