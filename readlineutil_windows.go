// +build windows

package readlineutil

import (
	"errors"
	"io"
)

type Term struct {
	io.Writer
}

type Option func()

var ErrNotSupported = errors.New("not supported")

// NewTerm returns ErrNotSupported, as on Windows, for
// command history we rely on the one built into CMD.COM.
func NewTerm(historyFilename string) (*Term, error) {
	return nil, ErrNotSupported
}

func (t *Term) Close() error {
	return nil
}

func (t *Term) WritePrompt(_ string) error {
	return nil
}

func (t *Term) Scan() bool {
	return false
}

func (t *Term) Err() error {
	return io.EOF
}

func (t *Term) Text() string {
	return ""
}

func WithHistoryFile(_ string) Option {
	return func() {}
}
