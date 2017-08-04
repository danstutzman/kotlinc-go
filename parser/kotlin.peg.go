package parser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Fun struct {
	name string
	args []Arg
	body []Call
}
type Arg struct {
	name  string
	type_ string
}
type Call struct {
	name string
	args []StringLiteral
}
type StringLiteral string

func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}

var g = &grammar{
	rules: []*rule{
		{
			name: "Input",
			pos:  position{line: 27, col: 1, offset: 373},
			expr: &actionExpr{
				pos: position{line: 27, col: 10, offset: 382},
				run: (*parser).callonInput1,
				expr: &labeledExpr{
					pos:   position{line: 27, col: 10, offset: 382},
					label: "fun",
					expr: &ruleRefExpr{
						pos:  position{line: 27, col: 14, offset: 386},
						name: "Fun",
					},
				},
			},
		},
		{
			name: "Fun",
			pos:  position{line: 31, col: 1, offset: 413},
			expr: &actionExpr{
				pos: position{line: 31, col: 8, offset: 420},
				run: (*parser).callonFun1,
				expr: &seqExpr{
					pos: position{line: 31, col: 8, offset: 420},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 31, col: 8, offset: 420},
							val:        "fun",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 31, col: 14, offset: 426},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 31, col: 16, offset: 428},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 31, col: 21, offset: 433},
								name: "Ident",
							},
						},
						&litMatcher{
							pos:        position{line: 31, col: 27, offset: 439},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 31, col: 31, offset: 443},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 31, col: 33, offset: 445},
							label: "funArgs",
							expr: &zeroOrOneExpr{
								pos: position{line: 31, col: 41, offset: 453},
								expr: &ruleRefExpr{
									pos:  position{line: 31, col: 41, offset: 453},
									name: "FunArgs",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 31, col: 50, offset: 462},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 31, col: 52, offset: 464},
							val:        ")",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 31, col: 56, offset: 468},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 31, col: 58, offset: 470},
							label: "funBody",
							expr: &ruleRefExpr{
								pos:  position{line: 31, col: 66, offset: 478},
								name: "FunBody",
							},
						},
					},
				},
			},
		},
		{
			name: "FunArgs",
			pos:  position{line: 46, col: 1, offset: 684},
			expr: &actionExpr{
				pos: position{line: 46, col: 12, offset: 695},
				run: (*parser).callonFunArgs1,
				expr: &seqExpr{
					pos: position{line: 46, col: 12, offset: 695},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 46, col: 12, offset: 695},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 46, col: 18, offset: 701},
								name: "FunArg",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 46, col: 25, offset: 708},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 46, col: 27, offset: 710},
							label: "rest",
							expr: &zeroOrOneExpr{
								pos: position{line: 46, col: 32, offset: 715},
								expr: &seqExpr{
									pos: position{line: 46, col: 34, offset: 717},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 46, col: 34, offset: 717},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 46, col: 38, offset: 721},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 46, col: 40, offset: 723},
											name: "FunArgs",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "FunArg",
			pos:  position{line: 54, col: 1, offset: 870},
			expr: &actionExpr{
				pos: position{line: 54, col: 11, offset: 880},
				run: (*parser).callonFunArg1,
				expr: &seqExpr{
					pos: position{line: 54, col: 11, offset: 880},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 54, col: 11, offset: 880},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 54, col: 16, offset: 885},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 54, col: 22, offset: 891},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 54, col: 24, offset: 893},
							val:        ":",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 54, col: 28, offset: 897},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 54, col: 30, offset: 899},
							label: "type_",
							expr: &ruleRefExpr{
								pos:  position{line: 54, col: 36, offset: 905},
								name: "Type",
							},
						},
					},
				},
			},
		},
		{
			name: "FunBody",
			pos:  position{line: 61, col: 1, offset: 990},
			expr: &actionExpr{
				pos: position{line: 61, col: 12, offset: 1001},
				run: (*parser).callonFunBody1,
				expr: &seqExpr{
					pos: position{line: 61, col: 12, offset: 1001},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 61, col: 12, offset: 1001},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 61, col: 16, offset: 1005},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 61, col: 18, offset: 1007},
							label: "statements",
							expr: &zeroOrOneExpr{
								pos: position{line: 61, col: 29, offset: 1018},
								expr: &ruleRefExpr{
									pos:  position{line: 61, col: 29, offset: 1018},
									name: "Statements",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 61, col: 41, offset: 1030},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 61, col: 43, offset: 1032},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Statements",
			pos:  position{line: 69, col: 1, offset: 1144},
			expr: &actionExpr{
				pos: position{line: 69, col: 15, offset: 1158},
				run: (*parser).callonStatements1,
				expr: &seqExpr{
					pos: position{line: 69, col: 15, offset: 1158},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 69, col: 15, offset: 1158},
							label: "statement",
							expr: &ruleRefExpr{
								pos:  position{line: 69, col: 25, offset: 1168},
								name: "Statement",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 69, col: 35, offset: 1178},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 69, col: 37, offset: 1180},
							label: "statements",
							expr: &zeroOrOneExpr{
								pos: position{line: 69, col: 48, offset: 1191},
								expr: &seqExpr{
									pos: position{line: 69, col: 49, offset: 1192},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 69, col: 49, offset: 1192},
											val:        ";",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 69, col: 53, offset: 1196},
											name: "_",
										},
										&labeledExpr{
											pos:   position{line: 69, col: 55, offset: 1198},
											label: "statements",
											expr: &ruleRefExpr{
												pos:  position{line: 69, col: 66, offset: 1209},
												name: "Statements",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Statement",
			pos:  position{line: 77, col: 1, offset: 1375},
			expr: &actionExpr{
				pos: position{line: 77, col: 14, offset: 1388},
				run: (*parser).callonStatement1,
				expr: &seqExpr{
					pos: position{line: 77, col: 14, offset: 1388},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 77, col: 14, offset: 1388},
							label: "ident",
							expr: &ruleRefExpr{
								pos:  position{line: 77, col: 20, offset: 1394},
								name: "Ident",
							},
						},
						&litMatcher{
							pos:        position{line: 77, col: 26, offset: 1400},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 77, col: 30, offset: 1404},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 77, col: 32, offset: 1406},
							label: "callArgs",
							expr: &zeroOrOneExpr{
								pos: position{line: 77, col: 41, offset: 1415},
								expr: &ruleRefExpr{
									pos:  position{line: 77, col: 41, offset: 1415},
									name: "CallArgs",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 77, col: 51, offset: 1425},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 77, col: 53, offset: 1427},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "CallArgs",
			pos:  position{line: 90, col: 1, offset: 1626},
			expr: &actionExpr{
				pos: position{line: 90, col: 13, offset: 1638},
				run: (*parser).callonCallArgs1,
				expr: &seqExpr{
					pos: position{line: 90, col: 13, offset: 1638},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 90, col: 13, offset: 1638},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 90, col: 19, offset: 1644},
								name: "CallArg",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 90, col: 27, offset: 1652},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 90, col: 29, offset: 1654},
							label: "rest",
							expr: &zeroOrOneExpr{
								pos: position{line: 90, col: 34, offset: 1659},
								expr: &seqExpr{
									pos: position{line: 90, col: 36, offset: 1661},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 90, col: 36, offset: 1661},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 90, col: 40, offset: 1665},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 90, col: 42, offset: 1667},
											name: "CallArgs",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "CallArg",
			pos:  position{line: 98, col: 1, offset: 1849},
			expr: &actionExpr{
				pos: position{line: 98, col: 12, offset: 1860},
				run: (*parser).callonCallArg1,
				expr: &seqExpr{
					pos: position{line: 98, col: 12, offset: 1860},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 98, col: 12, offset: 1860},
							val:        "\"",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 98, col: 16, offset: 1864},
							label: "contents",
							expr: &zeroOrMoreExpr{
								pos: position{line: 98, col: 25, offset: 1873},
								expr: &charClassMatcher{
									pos:        position{line: 98, col: 25, offset: 1873},
									val:        "[^\"]",
									chars:      []rune{'"'},
									ignoreCase: false,
									inverted:   true,
								},
							},
						},
						&litMatcher{
							pos:        position{line: 98, col: 31, offset: 1879},
							val:        "\"",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Type",
			pos:  position{line: 102, col: 1, offset: 1924},
			expr: &actionExpr{
				pos: position{line: 102, col: 9, offset: 1932},
				run: (*parser).callonType1,
				expr: &oneOrMoreExpr{
					pos: position{line: 102, col: 9, offset: 1932},
					expr: &charClassMatcher{
						pos:        position{line: 102, col: 9, offset: 1932},
						val:        "[a-zA-Z0-9<>]",
						chars:      []rune{'<', '>'},
						ranges:     []rune{'a', 'z', 'A', 'Z', '0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "Ident",
			pos:  position{line: 106, col: 1, offset: 1981},
			expr: &actionExpr{
				pos: position{line: 106, col: 10, offset: 1990},
				run: (*parser).callonIdent1,
				expr: &oneOrMoreExpr{
					pos: position{line: 106, col: 10, offset: 1990},
					expr: &charClassMatcher{
						pos:        position{line: 106, col: 10, offset: 1990},
						val:        "[a-zA-Z0-9]",
						ranges:     []rune{'a', 'z', 'A', 'Z', '0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "Integer",
			pos:  position{line: 110, col: 1, offset: 2037},
			expr: &actionExpr{
				pos: position{line: 110, col: 12, offset: 2048},
				run: (*parser).callonInteger1,
				expr: &seqExpr{
					pos: position{line: 110, col: 12, offset: 2048},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 110, col: 12, offset: 2048},
							expr: &litMatcher{
								pos:        position{line: 110, col: 12, offset: 2048},
								val:        "-",
								ignoreCase: false,
							},
						},
						&oneOrMoreExpr{
							pos: position{line: 110, col: 17, offset: 2053},
							expr: &charClassMatcher{
								pos:        position{line: 110, col: 17, offset: 2053},
								val:        "[0-9]",
								ranges:     []rune{'0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
				},
			},
		},
		{
			name:        "_",
			displayName: "\"whitespace\"",
			pos:         position{line: 114, col: 1, offset: 2103},
			expr: &zeroOrMoreExpr{
				pos: position{line: 114, col: 19, offset: 2121},
				expr: &charClassMatcher{
					pos:        position{line: 114, col: 19, offset: 2121},
					val:        "[ \\n\\t\\r]",
					chars:      []rune{' ', '\n', '\t', '\r'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 116, col: 1, offset: 2133},
			expr: &notExpr{
				pos: position{line: 116, col: 8, offset: 2140},
				expr: &anyMatcher{
					line: 116, col: 9, offset: 2141,
				},
			},
		},
	},
}

func (c *current) onInput1(fun interface{}) (interface{}, error) {
	return fun, nil
}

func (p *parser) callonInput1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInput1(stack["fun"])
}

func (c *current) onFun1(name, funArgs, funBody interface{}) (interface{}, error) {
	var args []Arg
	if funArgs == nil {
		args = []Arg{}
	} else {
		args = funArgs.([]Arg)
	}

	return Fun{
		name: name.(string),
		args: args,
		body: funBody.([]Call),
	}, nil
}

func (p *parser) callonFun1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFun1(stack["name"], stack["funArgs"], stack["funBody"])
}

func (c *current) onFunArgs1(first, rest interface{}) (interface{}, error) {
	args := []Arg{first.(Arg)}
	if rest != nil {
		args = append(args, rest.([]interface{})[2].([]Arg)...)
	}
	return args, nil
}

func (p *parser) callonFunArgs1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFunArgs1(stack["first"], stack["rest"])
}

func (c *current) onFunArg1(name, type_ interface{}) (interface{}, error) {
	return Arg{
		name:  name.(string),
		type_: type_.(string),
	}, nil
}

func (p *parser) callonFunArg1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFunArg1(stack["name"], stack["type_"])
}

func (c *current) onFunBody1(statements interface{}) (interface{}, error) {
	if statements == nil {
		return []string{}, nil
	} else {
		return statements.([]Call), nil
	}
}

func (p *parser) callonFunBody1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFunBody1(stack["statements"])
}

func (c *current) onStatements1(statement, statements interface{}) (interface{}, error) {
	out := []Call{statement.(Call)}
	if statements != nil {
		out = append(out, statements.([]interface{})[2].([]Call)...)
	}
	return out, nil
}

func (p *parser) callonStatements1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStatements1(stack["statement"], stack["statements"])
}

func (c *current) onStatement1(ident, callArgs interface{}) (interface{}, error) {
	if callArgs == nil {
		return Call{
			name: ident.(string),
		}, nil
	} else {
		return Call{
			name: ident.(string),
			args: callArgs.([]StringLiteral),
		}, nil
	}
}

func (p *parser) callonStatement1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStatement1(stack["ident"], stack["callArgs"])
}

func (c *current) onCallArgs1(first, rest interface{}) (interface{}, error) {
	exprs := []StringLiteral{first.(StringLiteral)}
	if rest != nil {
		exprs = append(exprs, rest.([]interface{})[2].([]StringLiteral)...)
	}
	return exprs, nil
}

func (p *parser) callonCallArgs1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCallArgs1(stack["first"], stack["rest"])
}

func (c *current) onCallArg1(contents interface{}) (interface{}, error) {
	return StringLiteral(c.text), nil
}

func (p *parser) callonCallArg1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCallArg1(stack["contents"])
}

func (c *current) onType1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onType1()
}

func (c *current) onIdent1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonIdent1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdent1()
}

func (c *current) onInteger1() (interface{}, error) {
	return strconv.Atoi(string(c.text))
}

func (p *parser) callonInteger1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInteger1()
}

var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// GlobalStore creates an Option to set a key to a certain value in
// the globalStore.
func GlobalStore(key string, value interface{}) Option {
	return func(p *parser) Option {
		old := p.cur.globalStore[key]
		p.cur.globalStore[key] = value
		return GlobalStore(key, old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (i interface{}, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			err = closeErr
		}
	}()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match

	// the globalStore allows the parser to store arbitrary values
	globalStore map[string]interface{}
}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos  position
	expr interface{}
	run  func(*parser) (interface{}, error)
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos             position
	val             string
	basicLatinChars [128]bool
	chars           []rune
	ranges          []rune
	classes         []*unicode.RangeTable
	ignoreCase      bool
	inverted        bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner    error
	pos      position
	prefix   string
	expected []string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	p := &parser{
		filename: filename,
		errs:     new(errList),
		data:     b,
		pt:       savepoint{position: position{line: 1}},
		recover:  true,
		cur: current{
			globalStore: make(map[string]interface{}),
		},
		maxFailPos:      position{col: 1, line: 1},
		maxFailExpected: make([]string, 0, 20),
	}
	p.setOptions(opts)
	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v   interface{}
	b   bool
	end savepoint
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	depth   int
	recover bool
	debug   bool

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// stats
	exprCnt int

	// parse fail
	maxFailPos            position
	maxFailExpected       []string
	maxFailInvertExpected bool
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position, []string{})
}

func (p *parser) addErrAt(err error, pos position, expected []string) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, pos: pos, prefix: buf.String(), expected: expected}
	p.errs.add(pe)
}

func (p *parser) failAt(fail bool, pos position, want string) {
	// process fail if parsing fails and not inverted or parsing succeeds and invert is set
	if fail == p.maxFailInvertExpected {
		if pos.offset < p.maxFailPos.offset {
			return
		}

		if pos.offset > p.maxFailPos.offset {
			p.maxFailPos = pos
			p.maxFailExpected = p.maxFailExpected[:0]
		}

		if p.maxFailInvertExpected {
			want = "!" + want
		}
		p.maxFailExpected = append(p.maxFailExpected, want)
	}
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError {
		if n == 1 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	// start rule is rule [0]
	p.read() // advance to first rune
	val, ok := p.parseRule(g.rules[0])
	if !ok {
		if len(*p.errs) == 0 {
			// If parsing fails, but no errors have been recorded, the expected values
			// for the farthest parser position are returned as error.
			maxFailExpectedMap := make(map[string]struct{}, len(p.maxFailExpected))
			for _, v := range p.maxFailExpected {
				maxFailExpectedMap[v] = struct{}{}
			}
			expected := make([]string, 0, len(maxFailExpectedMap))
			eof := false
			if _, ok := maxFailExpectedMap["!."]; ok {
				delete(maxFailExpectedMap, "!.")
				eof = true
			}
			for k := range maxFailExpectedMap {
				expected = append(expected, k)
			}
			sort.Strings(expected)
			if eof {
				expected = append(expected, "EOF")
			}
			p.addErrAt(errors.New("no match found, expected: "+listJoin(expected, ", ", "or")), p.maxFailPos, expected)
		}
		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func listJoin(list []string, sep string, lastSep string) string {
	switch len(list) {
	case 0:
		return ""
	case 1:
		return list[0]
	default:
		return fmt.Sprintf("%s %s %s", strings.Join(list[:len(list)-1], sep), lastSep, list[len(list)-1])
	}
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.exprCnt++
	var val interface{}
	var ok bool
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position, []string{})
		}
		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restore(pt)
	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn != utf8.RuneError {
		start := p.pt
		p.read()
		p.failAt(true, start.position, ".")
		return p.sliceFrom(start), true
	}
	p.failAt(false, p.pt.position, ".")
	return nil, false
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	start := p.pt

	// can't match EOF
	if cur == utf8.RuneError {
		p.failAt(false, start.position, chr.val)
		return nil, false
	}

	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		p.failAt(true, start.position, chr.val)
		return p.sliceFrom(start), true
	}
	p.failAt(false, start.position, chr.val)
	return nil, false
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for _, alt := range ch.alternatives {
		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			return val, ok
		}
	}
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	ignoreCase := ""
	if lit.ignoreCase {
		ignoreCase = "i"
	}
	val := fmt.Sprintf("%q%s", lit.val, ignoreCase)
	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.failAt(false, start.position, val)
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	p.failAt(true, start.position, val)
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	p.pushV()
	p.maxFailInvertExpected = !p.maxFailInvertExpected
	_, ok := p.parseExpr(not.expr)
	p.maxFailInvertExpected = !p.maxFailInvertExpected
	p.popV()
	p.restore(pt)
	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	vals := make([]interface{}, 0, len(seq.exprs))

	pt := p.pt
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}
