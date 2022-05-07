package window

import (
	"regexp"
	"strings"
)

type Parser struct {
	text string
	pos  int
}

func StartParsing(text string) *Parser {
	return &Parser{
		text: strings.ToLower(text),
	}
}

// expectAny compares the remaining string against any of given alternatives
// returns the consumed alternative (empty string if not matched)
func (p *Parser) expectAny(alts []string) string {
	for _, alt := range alts {
		if p.expect(alt) {
			return alt
		}
	}

	return "" // matched nothing
}

// expect consumes the given string if it matched the current remainder
func (p *Parser) expect(expected string) bool {
	if len(expected) <= len(p.text[p.pos:]) && expected == p.text[p.pos:p.pos+len(expected)] {
		p.pos += len(expected)
		return true
	}
	return false
}

// consumeUntil consumes all bytes from the text until it meets one of the alternatives (the alt is NOT consumed)
// it returns the consumed string and the found alternative
// it advances the position to after the found alternative
func (p *Parser) consumeUntil(alts []string) (consumed, matchedAlt string) {
	for p.pos < len(p.text) {
		for _, alt := range alts {
			if len(alt) <= len(p.text[p.pos:]) && alt == p.text[p.pos:p.pos+len(alt)] {
				matchedAlt = alt
				return
			}
		}

		// not matched, consume the char
		consumed = string(append([]byte(consumed), p.text[p.pos]))
		p.pos += 1
	}
	return
}

// consumeRE consumes all characters that match the given pattern.
// Pattern matches only the beginning of the string
func (p *Parser) consumeRE(pattern string) (matched string) {
	if p.isEof() {
		return ""
	}
	re := regexp.MustCompile(pattern)
	matched = re.FindString(p.getRemainder())

	p.pos += len(matched)
	return matched
}

// eatWs advances the pos to the next non-whitespace character in the text
func (p *Parser) eatWs() (eaten int) {
	for !p.isEof() {
		ch := p.text[p.pos]
		if ch == ' ' || ch == '\t' || ch == '\n' {
			p.pos++
			eaten++
			continue
		}
		break
	}
	return
}

func (p *Parser) getRemainder() string { return p.text[p.pos:] }

// rollback reset the pos back to i steps
func (p *Parser) rollback(i int) { p.pos -= i }

func (p *Parser) rollbackAt(pos int) { p.pos = pos }

func (p *Parser) isEof() bool { return p.pos >= len(p.text) }
