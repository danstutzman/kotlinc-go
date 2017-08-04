package main

import (
	"fmt"
)

const (
	ILLEGAL Token = iota
	EOF
	WS

	// Literals
	IDENT

	// Misc characters
	LPAREN
	RPAREN
	COLON
	LANGLE
	RANGLE
	LCURLY
	RCURLY

	// Keywords
	FUN

	STRING_LITERAL
)

type Token uint8

func (token Token) toString() string {
	switch token {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case WS:
		return "WS"
	case IDENT:
		return "IDENT"
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case COLON:
		return "COLON"
	case LANGLE:
		return "LANGLE"
	case RANGLE:
		return "RANGLE"
	case LCURLY:
		return "LCURLY"
	case RCURLY:
		return "RCURLY"
	case FUN:
		return "FUN"
	case STRING_LITERAL:
		return "STRING_LITERAL"
	default:
		panic(fmt.Errorf("Unknown token %d", token))
	}
}
