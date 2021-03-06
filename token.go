package goexp

type TokenType int

const (
	Unknown TokenType = iota
	EOF

	LeftParen    // (
	RightParen   // )
	LeftBracket  // [
	RightBracket // ]
	LeftBrace    // {
	RightBrace   // }
	Comma        // ,
	Period       // .
	Add          // +
	Sub          // -
	Mul          // *
	Div          // /
	Modulo       // %
	Power        // **

	And // &&
	Or  // ||
	Not // !
	Xor // xor

	NotEqual     // !=
	Equal        // ==
	Greater      // >
	GreaterEqual // >=
	Less         // <
	LessEqual    // <=

	Identifier // main
	String     // "abc"
	Integer    // 123
	Float      // 12.34

	True
	False
	Nil
)

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Pos     int
}

// func (t Token) String() string {
// 	return fmt.Sprintf("%+v", t)
// }
