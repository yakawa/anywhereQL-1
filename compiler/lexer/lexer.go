package lexer

import (
	"github.com/anywhereQL/anywhereQL/common/helper"
	"github.com/anywhereQL/anywhereQL/common/token"
	"github.com/anywhereQL/anywhereQL/common/value"
)

type lexer struct {
	src        []rune
	currentPos int
	readPos    int
}

func new(src string) *lexer {
	l := &lexer{
		currentPos: 0,
		readPos:    0,
		src:        []rune(src),
	}
	l.readChar()
	return l
}

func Lex(s string) token.Tokens {
	l := new(s)
	tokens := l.tokenize()

	return tokens
}

func (l *lexer) getCurrentChar() rune {
	if l.currentPos >= len(l.src) {
		return 0
	}
	return l.src[l.currentPos]
}

func (l *lexer) getNextChar() rune {
	if l.currentPos+1 >= len(l.src) {
		return 0
	}
	return l.src[l.currentPos+1]
}

func (l *lexer) readChar() {
	l.currentPos = l.readPos
	l.readPos++
}

func (l *lexer) tokenize() token.Tokens {
	tokens := token.Tokens{}

	for {
		ch := l.getCurrentChar()
		if ch == 0 {
			break
		} else if helper.IsWhiteSpace(ch) {
			l.readChar()
		} else {
			t, err := l.findToken()
			if err != nil {
				break
			}
			tokens = append(tokens, t)
		}
	}

	tokens = append(tokens, token.Token{Type: token.EOS})

	return tokens
}

func (l *lexer) findToken() (token.Token, error) {
	ch := l.getCurrentChar()
	if helper.IsSymbol(ch) {
		switch ch {
		case '"', '\'':
			v := l.readString(ch)
			val := value.Value{Type: value.STRING, String: v}
			return token.Token{Type: token.STRING, Literal: v, Value: val}, nil
		default:
			v, t := l.lookupSymbol()
			return token.Token{Type: t, Literal: v}, nil
		}
	} else if helper.IsDigit(ch) || (ch == '.' && helper.IsDigit(l.getNextChar())) {
		v := l.readNumber()
		val, err := value.Convert(v)
		if err != nil {
			return token.Token{}, err
		}
		t := token.Token{
			Type:    token.NUMBER,
			Literal: v,
			Value:   val,
		}
		return t, nil
	} else {
		v := l.readIdent()
		isKeyword, t := token.LookupKeyword(v)
		if isKeyword {
			return token.Token{
				Type:    t,
				Literal: v,
			}, nil
		} else {
			return token.Token{
				Type:    token.IDENT,
				Literal: v,
			}, nil
		}
	}
}

func (l *lexer) readString(quote rune) string {
	v := []rune("")
	l.readChar()
	for {
		ch := l.getCurrentChar()
		if ch == 0 || ch == '\n' {
			break
		}
		if ch == quote {
			l.readChar()
			if l.getCurrentChar() != quote {
				break
			}
		}
		v = append(v, ch)
		l.readChar()
	}
	return string(v)
}

func (l *lexer) readIdent() string {
	v := []rune("")
	for {
		ch := l.getCurrentChar()
		if ch == 0 || helper.IsWhiteSpace(ch) || helper.IsSymbol(ch) {
			break
		}
		v = append(v, ch)
		l.readChar()
	}
	return string(v)
}

func (l *lexer) readNumber() string {
	v := []rune("")
	for {
		ch := l.getCurrentChar()
		if helper.IsDigit(ch) {
			v = append(v, ch)
			l.readChar()
		} else {
			break
		}
	}
	ch := l.getCurrentChar()
	if ch == '.' {
		v = append(v, ch)
		l.readChar()
		for {
			ch = l.getCurrentChar()
			if helper.IsDigit(ch) {
				v = append(v, ch)
				l.readChar()
			} else {
				break
			}
		}
	}
	return string(v)
}

func (l *lexer) lookupSymbol() (string, token.Type) {
	ch := l.getCurrentChar()
	var t token.Type
	v := string(ch)
	switch ch {
	case ';':
		t = token.S_SEMICOLON
	case '+':
		t = token.S_PLUS
	case '-':
		t = token.S_MINUS
	case '*':
		t = token.S_ASTERISK
	case '/':
		t = token.S_SOLIDAS
	case '%':
		t = token.S_PERCENT
	case '(':
		t = token.S_LPAREN
	case ')':
		t = token.S_RPAREN
	case ',':
		t = token.S_COMMA
	case '"':
		t = token.S_DQUOTE
	case '\'':
		t = token.S_QUOTE
	case '.':
		t = token.S_PERIOD
	case '=':
		t = token.S_EQUAL
	case '<':
		if l.getNextChar() == '>' {
			l.readChar()
			t = token.S_NOT_EQUAL
			v = "<>"
		} else if l.getNextChar() == '=' {
			l.readChar()
			t = token.S_LESS_THAN_EQUAL
			v = "<="
		} else {
			t = token.S_LESS_THAN
		}
	case '>':
		if l.getNextChar() == '=' {
			l.readChar()
			t = token.S_GREATER_THAN_EQUAL
			v = ">="
		} else {
			t = token.S_GREATER_THAN
		}
	default:
		t = token.UNKNOWN
		v = ""
	}
	l.readChar()
	return v, t
}
