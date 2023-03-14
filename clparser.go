// clparser is a simple command-line parser.
package clparser

import (
	"fmt"
	"strings"
)

type CLParser struct {
	enableBackslashEscapes bool
}

func NewCLParser() *CLParser {
	return &CLParser{false}
}

func (clp *CLParser) BackslashEscapes(enableBackslashEscapes bool) *CLParser {
	clp.enableBackslashEscapes = enableBackslashEscapes
	return clp
}

const (
	stateSpace         = iota // white spaces
	statePlain                // normal characters
	stateEscape               // [\]
	stateSQuote               // [']
	stateDQuote               // ["]
	stateDQuotedEscape        // [\] in ["]
)

// Parse parse and split command line string.
func (clp *CLParser) Parse(commandLine string) ([]string, error) {
	var args []string = []string{}
	var arg strings.Builder
	arg.Grow(len(commandLine))
	state := stateSpace
	for _, ch := range commandLine {
		switch state {
		case stateSpace, statePlain:
			switch ch {
			case '\\':
				state = stateEscape
			case '\'':
				state = stateSQuote
			case '"':
				state = stateDQuote
			case ' ', '\t', '\r', '\n', '\f', '\v':
				if state == statePlain {
					args = append(args, arg.String())
					arg.Reset()
					state = stateSpace
				}
			default:
				arg.WriteRune(ch)
				state = statePlain
			}
		case stateEscape:
			arg.WriteRune(ch)
			state = statePlain
		case stateSQuote:
			switch ch {
			case '\'':
				state = statePlain
			default:
				arg.WriteRune(ch)
			}
		case stateDQuote:
			switch ch {
			case '\\':
				state = stateDQuotedEscape
			case '"':
				state = statePlain
			default:
				arg.WriteRune(ch)
			}
		case stateDQuotedEscape:
			if clp.enableBackslashEscapes {
				switch ch {
				case 'a':
					arg.WriteRune('\a')
				case 'b':
					arg.WriteRune('\b')
				case 'e', 'E':
					arg.WriteRune('\u001b')
				case 'f':
					arg.WriteRune('\f')
				case 'n':
					arg.WriteRune('\n')
				case 'r':
					arg.WriteRune('\r')
				case 't':
					arg.WriteRune('\t')
				case 'v':
					arg.WriteRune('\v')
				case '\\', '"':
					arg.WriteRune(ch)
				default:
					arg.WriteRune('\\')
					arg.WriteRune(ch)
				}
			} else {
				switch ch {
				case '\\', '"':
					arg.WriteRune(ch)
				default:
					arg.WriteRune('\\')
					arg.WriteRune(ch)
				}
			}
			state = stateDQuote
		}
	}
	switch state {
	case stateSpace:
		return args, nil
	case statePlain:
		args = append(args, arg.String())
		return args, nil
	case stateEscape:
		return nil, fmt.Errorf("terminated by an escape character: [%s]", commandLine)
	case stateSQuote:
		return nil, fmt.Errorf("single-quoted string not closed: [%s]", commandLine)
	case stateDQuote, stateDQuotedEscape:
		return nil, fmt.Errorf("double-quoted string not closed: [%s]", commandLine)
	default:
		return nil, fmt.Errorf("unexpected error: [%s]", commandLine)
	}
}
