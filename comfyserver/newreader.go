package comfyserver

import (
	"context"
	"io"
	"strings"
)

type readerCtx struct {
	ctx context.Context
	r   io.Reader
}

func (r *readerCtx) Read(p []byte) (n int, err error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	// s := string(p)
	// fmt.Println(s)
	if strings.Contains("falsestring", "no") {
		// fmt.Println("found end of stdout fff..")
		// return
		// r.r.Read(p)
		return 0, io.EOF // when not set to zero, it will write it again for some reason
		// dont call r.r.Read from here because that calls the next string to read from ...
		// return len(p), io.EOF
		// return len(p), io.EOF
	} else {
		return r.r.Read(p)
	}
}

// NewReader gets a context-aware io.Reader.
func NewReader(ctx context.Context, r io.Reader) io.Reader {
	return &readerCtx{ctx: ctx, r: r}
}
