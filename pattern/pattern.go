package pattern

import (
	"fmt"
	"regexp"
	"strings"
)

// Type describes a type of pattern.
type Type int

// Pattern types.
const (
	Text Type = iota
	Regexp
)

// Pattern type prefixes.
const (
	PrefixRegexp = "re"
)

// Pattern matches string patterns by performing an exact text match or a
// regular expression match.
type Pattern struct {
	value string
	t     Type
	re    *regexp.Regexp // Compiled expression when value is in regexp form
}

// String returns the string form of the pattern, which will include a prefix
// for non-text patterns.
func (p *Pattern) String() string {
	return p.value
}

// Type returns the type of the pattern.
func (p *Pattern) Type() Type {
	return p.t
}

// Match returns true if the pattern matches the string.
//
// If the pattern is a regular expression with one or more capturing subgroups,
// the non-empty submatches will be returned as well.
func (p *Pattern) Match(s string) (submatches []string, matched bool) {
	switch p.t {
	case Text:
		return nil, p.value == s
	case Regexp:
		if p.re.NumSubexp() == 0 {
			return nil, p.re.MatchString(s)
		}
		subm := p.re.FindStringSubmatch(s)
		if subm == nil {
			return nil, false // Failed match
		}
		if len(subm) < 2 {
			return nil, true // Successful match with non-capturing subexpressions
		}
		submatches = make([]string, 0, len(subm)-1)
		for _, submatch := range subm {
			if submatch != "" {
				submatches = append(submatches, submatch)
			}
		}
		return submatches, true
	default:
		return nil, false
	}
}

// MarshalText saves the string-encoded form of the pattern in text.
func (p *Pattern) MarshalText() (text []byte, err error) {
	return []byte(p.value), nil
}

// UnmarshalText parses the given text and stores the parsed result in p.
//
// If the pattern describes a regular expression and that expression cannot be
// compiled, an error will be returned.
func (p *Pattern) UnmarshalText(text []byte) error {
	p.t = Text
	p.value = string(text)
	if parts := strings.SplitN(p.value, ":", 2); len(parts) == 2 {
		switch parts[0] {
		case PrefixRegexp:
			var err error
			p.t = Regexp
			p.re, err = regexp.Compile(parts[1])
			if err != nil {
				return fmt.Errorf("pattern regular expression compilation error: %s", err)
			}
		}
	}
	return nil
}
