//go:generate stringer -type=TokenType -output=lex_string.go

package tis100mc

import "fmt"

const eof = -1

type TokenType int

const (
	TokenComment TokenType = iota
	TokenLabel
	TokenInstruction
	TokenArgument
	TokenComma
)

type stateFn func(*lexer) stateFn

type Token struct {
	typ        TokenType
	start, end int

	// TODO remove
	val string
}

func (t Token) String() string {
	return fmt.Sprintf("%s - %s", t.typ, t.val)
}

type TokenReceiver interface {
	Receive(token Token)
}

type lexer struct {
	// TODO make this an io.Reader with intent to support
	// never ending streams.
	bytes []byte

	receiver TokenReceiver
	state    stateFn
	pos      int
	start    int
}

func NewLexer(receiver TokenReceiver) *lexer {
	if receiver == nil {
		panic("lexer does not accept nil TokenReceiver.")
	}
	return &lexer{
		state:    lexWhitespace,
		receiver: receiver,
	}
}

func (l *lexer) Run() {
	for l.state != nil {
		l.state = l.state(l)
	}
}

// TODO should reconsider how tokens are emitted. lexer could feed over
// a channel or otherwise.
func (l *lexer) emit(t TokenType) {
	l.receiver.Receive(Token{
		typ:   t,
		val:   string(l.bytes[l.start:l.pos]),
		start: l.start,
		end:   l.pos,
	})
}

func (l *lexer) next() {
	l.pos++
}

func (l *lexer) reset() {
	l.start = l.pos
}

func (l *lexer) discard() {
	l.pos++
	l.start = l.pos
}

func (l *lexer) rune() rune {
	if l.pos >= len(l.bytes) {
		return eof
	}
	return rune(l.bytes[l.pos])
}

func lexWhitespace(l *lexer) stateFn {
	for {
		switch l.rune() {
		case ' ', '\t', '\n':
			l.discard()
		case eof:
			return nil
		default:
			return lexInstruction
		}
	}
}

func lexInstruction(l *lexer) stateFn {
	for {
		switch l.rune() {
		case '#':
			if l.pos != l.start {
				l.emit(TokenInstruction)
				l.reset()
			}
			return lexComment
		case ' ':
			l.emit(TokenInstruction)
			l.discard()
			return lexArgument
		case '\n':
			l.emit(TokenInstruction)
			l.discard()
			return lexWhitespace
		case ':':
			l.emit(TokenLabel)
			l.discard()
			return lexWhitespace
		case eof:
			l.emit(TokenLabel)
			return nil
		default:
			l.next()
		}
	}
}

func lexArgument(l *lexer) stateFn {
	for {
		switch l.rune() {
		case ',':
			if l.pos != l.start {
				l.emit(TokenArgument)
			}
			l.reset()
			l.emit(TokenComma)
			l.discard()
		case '#':
			if l.pos != l.start {
				l.emit(TokenArgument)
				l.reset()
			}
			return lexComment
		case ' ', '\t':
			if l.pos != l.start {
				l.emit(TokenArgument)
			}
			l.discard()
		case '\n':
			if l.pos != l.start {
				l.emit(TokenArgument)
			}
			l.discard()
			return lexWhitespace
		case eof:
			l.emit(TokenArgument)
			return nil
		default:
			l.next()
		}
	}
}

func lexComment(l *lexer) stateFn {
	for {
		switch l.rune() {
		case '\n':
			l.emit(TokenComment)
			l.discard()
			return lexWhitespace
		case eof:
			return nil
		default:
			l.next()
		}
	}
}
