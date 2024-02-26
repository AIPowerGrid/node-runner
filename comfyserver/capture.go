package comfyserver

import (
	"bytes"
	"io"
	"strings"
)

type CapturingPassThroughWriter struct {
	buf          bytes.Buffer
	w            io.Writer
	initChan     chan string
	targetString string
}

// NewCapturingPassThroughWriter creates new CapturingPassThroughWriter
func NewCapturingPassThroughWriter(w io.Writer, targetString string, ch chan string) *CapturingPassThroughWriter {
	return &CapturingPassThroughWriter{
		w:            w,
		initChan:     ch,
		targetString: targetString,
	}
}

func (w *CapturingPassThroughWriter) Write(d []byte) (int, error) {
	n, err := w.buf.Write(d)
	n, err = w.w.Write(d)
	s := string(d)
	if strings.Contains(s, w.targetString) {
		w.initChan <- "found"
		close(w.initChan)
	}
	return n, err
}

// Bytes returns bytes written to the writer
func (w *CapturingPassThroughWriter) Bytes() []byte {
	return w.buf.Bytes()
}
