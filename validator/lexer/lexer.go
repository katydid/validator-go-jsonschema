// Code generated by gocc; DO NOT EDIT.

package lexer

import (
	"os"
	"unicode/utf8"

	"github.com/katydid/validator-go-jsonschema/validator/token"
)

const (
	NoState    = -1
	NumStates  = 277
	NumSymbols = 205
)

type Lexer struct {
	src     []byte
	pos     int
	line    int
	column  int
	Context token.Context
}

func NewLexer(src []byte) *Lexer {
	lexer := &Lexer{
		src:     src,
		pos:     0,
		line:    1,
		column:  1,
		Context: nil,
	}
	return lexer
}

// SourceContext is a simple instance of a token.Context which
// contains the name of the source file.
type SourceContext struct {
	Filepath string
}

func (s *SourceContext) Source() string {
	return s.Filepath
}

func NewLexerFile(fpath string) (*Lexer, error) {
	src, err := os.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	lexer := NewLexer(src)
	lexer.Context = &SourceContext{Filepath: fpath}
	return lexer, nil
}

func (l *Lexer) Scan() (tok *token.Token) {
	tok = &token.Token{}
	if l.pos >= len(l.src) {
		tok.Type = token.EOF
		tok.Pos.Offset, tok.Pos.Line, tok.Pos.Column = l.pos, l.line, l.column
		tok.Pos.Context = l.Context
		return
	}
	start, startLine, startColumn, end := l.pos, l.line, l.column, 0
	tok.Type = token.INVALID
	state, rune1, size := 0, rune(-1), 0
	for state != -1 {
		if l.pos >= len(l.src) {
			rune1 = -1
		} else {
			rune1, size = utf8.DecodeRune(l.src[l.pos:])
			l.pos += size
		}

		nextState := -1
		if rune1 != -1 {
			nextState = TransTab[state](rune1)
		}
		state = nextState

		if state != -1 {

			switch rune1 {
			case '\n':
				l.line++
				l.column = 1
			case '\r':
				l.column = 1
			case '\t':
				l.column += 4
			default:
				l.column++
			}

			switch {
			case ActTab[state].Accept != -1:
				tok.Type = ActTab[state].Accept
				end = l.pos
			case ActTab[state].Ignore != "":
				start, startLine, startColumn = l.pos, l.line, l.column
				state = 0
				if start >= len(l.src) {
					tok.Type = token.EOF
				}

			}
		} else {
			if tok.Type == token.INVALID {
				end = l.pos
			}
		}
	}
	if end > start {
		l.pos = end
		tok.Lit = l.src[start:end]
	} else {
		tok.Lit = []byte{}
	}
	tok.Pos.Offset, tok.Pos.Line, tok.Pos.Column = start, startLine, startColumn
	tok.Pos.Context = l.Context

	return
}

func (l *Lexer) Reset() {
	l.pos = 0
}

/*
Lexer symbols:
0: '('
1: ')'
2: '0'
3: '-'
4: '('
5: ')'
6: '('
7: '-'
8: ')'
9: '$'
10: '$'
11: '$'
12: '$'
13: '$'
14: '$'
15: '{'
16: ','
17: '}'
18: '<'
19: 'e'
20: 'm'
21: 'p'
22: 't'
23: 'y'
24: '>'
25: '['
26: ']'
27: 'b'
28: 'o'
29: 'o'
30: 'l'
31: '['
32: ']'
33: 'i'
34: 'n'
35: 't'
36: '['
37: ']'
38: 'u'
39: 'i'
40: 'n'
41: 't'
42: '['
43: ']'
44: 'd'
45: 'o'
46: 'u'
47: 'b'
48: 'l'
49: 'e'
50: '['
51: ']'
52: 's'
53: 't'
54: 'r'
55: 'i'
56: 'n'
57: 'g'
58: '['
59: ']'
60: '['
61: ']'
62: 'b'
63: 'y'
64: 't'
65: 'e'
66: 't'
67: 'r'
68: 'u'
69: 'e'
70: 'f'
71: 'a'
72: 'l'
73: 's'
74: 'e'
75: '='
76: '('
77: ')'
78: '{'
79: '}'
80: ','
81: ';'
82: '#'
83: '&'
84: '|'
85: '['
86: ']'
87: ':'
88: '!'
89: '*'
90: '_'
91: '.'
92: '@'
93: '-'
94: '>'
95: '='
96: '='
97: '!'
98: '='
99: '<'
100: '>'
101: '<'
102: '='
103: '>'
104: '='
105: '~'
106: '='
107: '*'
108: '='
109: '^'
110: '='
111: '$'
112: '='
113: ':'
114: ':'
115: '?'
116: '/'
117: '/'
118: '\n'
119: '/'
120: '*'
121: '*'
122: '*'
123: '/'
124: ' '
125: '\t'
126: '\n'
127: '\r'
128: '0'
129: '0'
130: 'x'
131: 'X'
132: '-'
133: 'e'
134: 'E'
135: '+'
136: '-'
137: '.'
138: '.'
139: '.'
140: '_'
141: '_'
142: 'd'
143: 'o'
144: 'u'
145: 'b'
146: 'l'
147: 'e'
148: 'i'
149: 'n'
150: 't'
151: 'u'
152: 'i'
153: 'n'
154: 't'
155: '['
156: ']'
157: 'b'
158: 'y'
159: 't'
160: 'e'
161: 's'
162: 't'
163: 'r'
164: 'i'
165: 'n'
166: 'g'
167: 'b'
168: 'o'
169: 'o'
170: 'l'
171: '.'
172: '\'
173: 'U'
174: '\'
175: 'u'
176: '\'
177: 'x'
178: '\'
179: '`'
180: '`'
181: '\'
182: 'a'
183: 'b'
184: 'f'
185: 'n'
186: 'r'
187: 't'
188: 'v'
189: '\'
190: '''
191: '"'
192: '"'
193: '"'
194: '''
195: '''
196: '0'-'9'
197: '0'-'7'
198: '0'-'9'
199: 'A'-'F'
200: 'a'-'f'
201: '1'-'9'
202: 'A'-'Z'
203: 'a'-'z'
204: .
*/
