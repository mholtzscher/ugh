package nlp

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

const tokenCapacityDivisor = 2

type tokenKind int

const (
	tokenWord tokenKind = iota
	tokenQuoted
	tokenColon
	tokenPlus
	tokenMinus
	tokenBang
	tokenLParen
	tokenRParen
	tokenEOF
)

type token struct {
	kind   tokenKind
	text   string
	quoted bool
	span   Span
}

type cursor struct {
	offset int
	line   int
	column int
}

func newCursor() cursor {
	return cursor{line: 1, column: 1}
}

func lex(input string) ([]token, error) {
	tokens := make([]token, 0, len(input)/tokenCapacityDivisor)
	cur := newCursor()

	for len(input) > 0 {
		r, width := utf8.DecodeRuneInString(input)
		if unicode.IsSpace(r) {
			cur = advanceCursor(cur, r, width)
			input = input[width:]
			continue
		}

		start := cur
		switch r {
		case ':':
			tokens = append(tokens, newSymbolToken(tokenColon, ":", start, r, width))
			cur = advanceCursor(cur, r, width)
			input = input[width:]
		case '+':
			tokens = append(tokens, newSymbolToken(tokenPlus, "+", start, r, width))
			cur = advanceCursor(cur, r, width)
			input = input[width:]
		case '-':
			tokens = append(tokens, newSymbolToken(tokenMinus, "-", start, r, width))
			cur = advanceCursor(cur, r, width)
			input = input[width:]
		case '!':
			tokens = append(tokens, newSymbolToken(tokenBang, "!", start, r, width))
			cur = advanceCursor(cur, r, width)
			input = input[width:]
		case '(':
			tokens = append(tokens, newSymbolToken(tokenLParen, "(", start, r, width))
			cur = advanceCursor(cur, r, width)
			input = input[width:]
		case ')':
			tokens = append(tokens, newSymbolToken(tokenRParen, ")", start, r, width))
			cur = advanceCursor(cur, r, width)
			input = input[width:]
		case '"':
			text, rest, end, err := lexQuoted(input, start)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, token{kind: tokenQuoted, text: text, quoted: true, span: spanFrom(start, end)})
			cur = end
			input = rest
		default:
			word, rest, end := lexWord(input, start)
			tokens = append(tokens, token{kind: tokenWord, text: word, span: spanFrom(start, end)})
			cur = end
			input = rest
		}
	}

	eofPos := positionFrom(cur)
	tokens = append(tokens, token{kind: tokenEOF, span: Span{Start: eofPos, End: eofPos}})
	return tokens, nil
}

func lexQuoted(input string, start cursor) (string, string, cursor, error) {
	cur := start
	_, width := utf8.DecodeRuneInString(input)
	cur = advanceCursor(cur, '"', width)
	input = input[width:]

	var out strings.Builder
	escaped := false

	for len(input) > 0 {
		r, w := utf8.DecodeRuneInString(input)
		if escaped {
			out.WriteRune(r)
			escaped = false
			cur = advanceCursor(cur, r, w)
			input = input[w:]
			continue
		}
		if r == '\\' {
			escaped = true
			cur = advanceCursor(cur, r, w)
			input = input[w:]
			continue
		}
		if r == '"' {
			cur = advanceCursor(cur, r, w)
			input = input[w:]
			return out.String(), input, cur, nil
		}
		out.WriteRune(r)
		cur = advanceCursor(cur, r, w)
		input = input[w:]
	}

	return "", "", cur, errors.New("unterminated quoted string")
}

func lexWord(input string, start cursor) (string, string, cursor) {
	cur := start
	var out strings.Builder
	for len(input) > 0 {
		r, width := utf8.DecodeRuneInString(input)
		if unicode.IsSpace(r) || isDelimiterRune(r) {
			break
		}
		out.WriteRune(r)
		cur = advanceCursor(cur, r, width)
		input = input[width:]
	}
	return out.String(), input, cur
}

func newSymbolToken(kind tokenKind, text string, start cursor, r rune, width int) token {
	return token{kind: kind, text: text, span: spanFrom(start, advanceCursor(start, r, width))}
}

func isDelimiterRune(r rune) bool {
	switch r {
	case ':', '+', '-', '!', '(', ')', '"':
		return true
	default:
		return false
	}
}

func advanceCursor(cur cursor, r rune, width int) cursor {
	cur.offset += width
	if r == '\n' {
		cur.line++
		cur.column = 1
		return cur
	}
	cur.column++
	return cur
}

func spanFrom(start cursor, end cursor) Span {
	return Span{Start: positionFrom(start), End: positionFrom(end)}
}

func positionFrom(cur cursor) Position {
	return Position{Offset: cur.offset, Line: cur.line, Column: cur.column}
}
