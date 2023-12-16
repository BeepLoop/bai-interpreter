package lexer

import "strconv"

type Scanner struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
}

func CreateScanner(source string) *Scanner {
	return &Scanner{
		source:  source,
		start:   0,
		current: 0,
		line:    1,
	}
}

func (S *Scanner) ScanTokens() []Token {
	for S.IsAtEnd() == false {
		S.start = S.current
		S.ScanToken()
	}

	S.tokens = append(S.tokens, *CreateToken(EOF, "", nil, S.line))
	return S.tokens
}

func (S *Scanner) ScanToken() {
	c := S.Advance()

	switch c {
	case '(':
		S.AddToken(LEFT_PAREN, nil)
	case ')':
		S.AddToken(RIGHT_PAREN, nil)
	case '{':
		S.AddToken(LEFT_BRACE, nil)
	case '}':
		S.AddToken(RIGHT_BRACE, nil)
	case ',':
		S.AddToken(COMMA, nil)
	case '.':
		S.AddToken(DOT, nil)
	case '-':
		S.AddToken(MINUS, nil)
	case '+':
		S.AddToken(PLUS, nil)
	case ';':
		S.AddToken(SEMICOLON, nil)
	case '*':
		S.AddToken(STAR, nil)
	case '!':
		if S.Match('=') {
			S.AddToken(BANG_EQUAL, nil)
		} else {
			S.AddToken(BANG, nil)
		}
	case '=':
		if S.Match('=') {
			S.AddToken(EQUAL_EQUAL, nil)
		} else {
			S.AddToken(EQUAL, nil)
		}
	case '<':
		if S.Match('=') {
			S.AddToken(LESS_EQUAL, nil)
		} else {
			S.AddToken(LESS, nil)
		}
	case '>':
		if S.Match('=') {
			S.AddToken(GREATER_EQUAL, nil)
		} else {
			S.AddToken(GREATER, nil)
		}
	case '/':
		if S.Match('/') {
			// detect a comment
			for S.Peek() != '\n' && S.IsAtEnd() == false {
				S.Advance()
			}
		} else {
			S.AddToken(SLASH, nil)
		}
	case '"':
		S.String()
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		S.line++
	default:
		if S.IsDigit(c) == true {
			S.Number()
		} else {
			Error(S.line, "Unexpected character.")
		}
	}
}

func (S *Scanner) Number() {
	for S.IsDigit(S.Peek()) {
		S.Advance()
	}

	if S.Peek() == '.' && S.IsDigit(S.PeekNext()) {
		S.Advance()

		for S.IsDigit(S.Peek()) {
			S.Advance()
		}
	}

	value, err := strconv.ParseFloat(S.source[S.start:S.current], 32)
	if err != nil {
		Error(S.line, "Cannot parse to Number")
	}

	S.AddToken(NUMBER, value)
}

func (S *Scanner) String() {
	for S.Peek() != '"' && S.IsAtEnd() == false {
		if S.Peek() == '\n' {
			S.line++
		}
		S.Advance()
	}

	if S.IsAtEnd() {
		Error(S.line, "Unterminated String")
		return
	}

	S.Advance()

	value := S.source[S.start+1 : S.current-1]
	S.AddToken(STRING, value)
}

func (S *Scanner) Match(expected byte) bool {
	if S.IsAtEnd() {
		return false
	}
	if S.source[S.current] != expected {
		return false
	}

	S.current++
	return true
}

func (S *Scanner) Peek() byte {
	if S.IsAtEnd() {
		return 0
	}
	return S.source[S.current]
}

func (S *Scanner) PeekNext() byte {
	if S.current+1 > len(S.source) {
		return 0
	}
	return S.source[S.current+1]
}

func (S *Scanner) IsDigit(c byte) bool {
	return c >= 48 && c <= 57
}

func (S *Scanner) Advance() byte {
	S.current++
	return S.source[S.current-1]
}

func (S *Scanner) AddToken(t_type TokenType, literal any) {
	text := S.source[S.start:S.current]
	S.tokens = append(S.tokens, *CreateToken(t_type, text, literal, S.line))
}

func (S *Scanner) IsAtEnd() bool {
	res := S.current >= len(S.source)
	return res
}
