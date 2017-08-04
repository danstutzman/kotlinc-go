package main

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type Scanner struct {
	r *bufio.Reader
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

func (s *Scanner) Scan() (token Token, matched string) {
	// Read the next rune.
	ch := s.read()

	if isWhitespace(ch) {
		// If we see whitespace then consume all contiguous whitespace.
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(ch) {
		// If we see a letter then consume as an ident or reserved word.
		s.unread()
		return s.scanIdent()
	} else if ch == '"' {
		return s.scanStringLiteral()
	}

	// Otherwise read the individual character.
	switch ch {
	case eof:
		return EOF, ""
	case '(':
		return LPAREN, string(ch)
	case ')':
		return RPAREN, string(ch)
	case ':':
		return COLON, string(ch)
	case '<':
		return LANGLE, string(ch)
	case '>':
		return RANGLE, string(ch)
	case '{':
		return LCURLY, string(ch)
	case '}':
		return RCURLY, string(ch)
	}

	return ILLEGAL, string(ch)
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanWhitespace() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

// scanIdent consumes the current rune and all contiguous ident runes.
func (s *Scanner) scanIdent() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	// If the string matches a keyword then return that keyword.
	switch strings.ToUpper(buf.String()) {
	case "FUN":
		return FUN, buf.String()
	}

	// Otherwise return as a regular identifier.
	return IDENT, buf.String()
}

func (s *Scanner) scanStringLiteral() (tok Token, lit string) {
	var buf bytes.Buffer

	for {
		if ch := s.read(); ch == eof {
			break
		} else if ch == '"' {
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return STRING_LITERAL, buf.String()
}

// read reads the next rune from the buffered reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() {
	_ = s.r.UnreadRune()
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

var eof = rune(0)