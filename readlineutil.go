// Package readlineutil provides a bufio.Scanner like interface to
// a GNU readline-style package, enabling command line programs
// to add a command history, browsable using arrow keys, to their interfaces.
//
// +build !windows

package readlineutil

import (
	"errors"
	"io"
	"os"

	"github.com/chzyer/readline"
)

type Term struct {
	inst *readline.Instance
	io.Writer
	line       string
	err        error
	prevPrompt string
}

type Option func(*conf)

type conf struct {
	rc *readline.Config
}

var ErrNotSupported = errors.New("not supported")

// NewTerm wraps a chzyer/readline.Instance so that it can be used
// in a way similar to bufio.Scanner. NewTerm does not check whether
// os.Stdin actually is a terminal. A nil Term and ErrNotSupported are
// returned in case NewTerm is called with TERM environment
// variable set to dumb.
func NewTerm(options ...Option) (*Term, error) {
	var c conf

	if os.Getenv("TERM") == "dumb" {
		return nil, ErrNotSupported
	}
	c.rc = new(readline.Config)
	c.rc.DisableAutoSaveHistory = true
	for _, o := range options {
		o(&c)
	}
	inst, err := readline.NewEx(c.rc)
	if err != nil {
		return nil, err
	}
	t := new(Term)
	t.inst = inst
	t.Writer = inst.Stdout()
	return t, nil
}

func (t *Term) Close() error {
	return t.inst.Close()
}

func (t *Term) WritePrompt(p string) error {
	if p != t.prevPrompt {
		t.inst.SetPrompt(p)
		t.prevPrompt = p
	}
	return nil
}

func (t *Term) Scan() bool {
	line, err := t.inst.Readline()
	if err != nil {
		t.err = err
		return false
	}
	t.line = line
	if t.prevPrompt != "" && line != "" {
		t.inst.SaveHistory(line)
	}
	return true
}

func (t *Term) Err() error {
	if t.err == io.EOF {
		return nil
	}
	return t.err
}

func (t *Term) Text() string {
	return t.line
}

// WithHistoryFile configures a filename to be used to persist
// the command history.
func WithHistoryFile(filename string) Option {
	return func(c *conf) {
		c.rc.HistoryFile = filename
	}
}
