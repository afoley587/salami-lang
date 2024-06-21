package lexer

import (
	"bufio"
	"io"
	"unicode"

	"github.com/afoley/salami-lang/tok"
)

type LexPosition struct {
	Line   int
	Column int
}

type Lexer struct {
	pos    LexPosition
	reader *bufio.Reader
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		pos:    LexPosition{Line: 1, Column: 0},
		reader: bufio.NewReader(reader),
	}
}

func (l *Lexer) Lex() (LexPosition, tok.TokenType, string) {
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return l.pos, tok.EOF, ""
			}

			panic(err)
		}

		l.pos.Column++
		switch r {
		case '\n':
			l.handleNewLine()
		case '=':
			return l.pos, tok.ASSIGN, "="
		case '+':
			return l.pos, tok.PLUS, "+"
		case '-':
			return l.pos, tok.MINUS, "-"
		case '*':
			return l.pos, tok.ASTERISK, "*"
		case '/':
			return l.pos, tok.SLASH, "/"
		case ';':
			return l.pos, tok.SEMICOLON, ";"
		case '(':
			return l.pos, tok.LPAREN, "("
		case ')':
			return l.pos, tok.RPAREN, ")"
		case '{':
			return l.pos, tok.LBRACE, "{"
		case '}':
			return l.pos, tok.RBRACE, "}"
		case '>':
			return l.pos, tok.GT, ">"
		case '<':
			return l.pos, tok.LT, "<"
		case ',':
			return l.pos, tok.COMMA, ","
		default:
			if unicode.IsSpace(r) {
				continue // nothing to do here, just move on
			} else if unicode.IsDigit(r) {
				starts := l.pos
				l.goBack()
				literal := l.readDigit()
				return starts, tok.INT, literal
			} else if unicode.IsPrint(r) {
				starts := l.pos
				l.goBack()
				literal := l.readIdentifier()
				return starts, tok.KeywordLookup(literal), literal

			} else {
				return l.pos, tok.ILLEGAL, string(r)
			}
		}
	}
}

func (l *Lexer) NextToken() tok.Tok {
	_, tokType, literal := l.Lex()
	return tok.Tok{Type: tokType, Literal: literal}
}

func (l *Lexer) handleNewLine() {
	l.pos.Line++
	l.pos.Column = 0
}

func (l *Lexer) goBack() {

	if err := l.reader.UnreadRune(); err != nil {
		panic(err)
	}

	l.pos.Column--
}

func (l *Lexer) readDigit() string {
	literal := ""

	for {
		r, _, err := l.reader.ReadRune()

		if err != nil {
			if err == io.EOF {
				return literal
			}
		}

		l.pos.Column++

		if unicode.IsDigit(r) {
			literal += string(r)
		} else {
			l.goBack()
			return literal
		}
	}
}

func (l *Lexer) readIdentifier() string {
	literal := ""

	for {
		r, _, err := l.reader.ReadRune()

		if err != nil {
			if err == io.EOF {
				return literal
			}
		}

		l.pos.Column++

		if unicode.IsLetter(r) {
			literal += string(r)
		} else {
			l.goBack()
			return literal
		}
	}
}
