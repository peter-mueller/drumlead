package parser

import (
	"fmt"
	"log"
	"strings"
	"unicode/utf8"
)

type itemType int

const (
	ERROR itemType = iota

	HASHTAG
	TEXT
)

const (
	eof        rune = -1
	hashtag         = '#'
	whitespace      = " \t"
	newline         = '\n'
)

type item struct {
	typ   itemType
	value string
}

func Lex(input string) <-chan item {
	c := make(chan item, 1000)

	go func() {
		defer close(c)
		l := &lexer{
			input: input,
			items: c,
		}

		for state := lexLine; state != nil; {
			state = state(l)
		}
	}()

	return c
}

type stateFunc func(*lexer) stateFunc

type lexer struct {
	input string

	start int
	pos   int
	width int

	items chan item
}

func lexLine(l *lexer) stateFunc {

	l.accept(whitespace)
	l.ignore()

	r := l.next()
	switch r {
	case hashtag:
		l.emit(HASHTAG)
		return lexLine
	case eof:
		return nil
	case newline:
		l.ignore()
		l.emit(TEXT)
		return lexLine
	default:
		return lexText
	}
}

func lexText(l *lexer) stateFunc {
	i := strings.IndexRune(l.current(), newline)
	if i == -1 {
		l.pos = len(l.input)
		l.emit(TEXT)
		return nil
	}
	l.pos += i

	l.emit(TEXT)
	l.must(newline)
	l.ignore()

	return lexLine
}

func (l *lexer) must(expected rune) {
	r := l.next()
	if r != expected {
		log.Fatalf("expected %s but got %s", string(expected), string(r))
	}
}

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) current() string {
	return l.input[l.pos:]
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) emit(t itemType) {
	l.items <- item{
		typ:   t,
		value: l.input[l.start:l.pos],
	}
	l.start = l.pos
}
func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}
func (l *lexer) errorf(format string, args ...interface{}) stateFunc {
	l.items <- item{
		typ:   ERROR,
		value: fmt.Sprintf(format, args...),
	}
	return nil
}
