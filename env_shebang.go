package enry

import "strings"

// splitEnvShebang handles env -S quoting and escapes without executing commands
// or expanding variables from the process environment.
func splitEnvShebang(source string) ([]string, bool) {
	var words []string
	var word strings.Builder
	var quote byte
	started := false
	flush := func() {
		if started {
			words = append(words, word.String())
			word.Reset()
			started = false
		}
	}
	for i := 0; i < len(source); i++ {
		c := source[i]
		if c == '\\' {
			if i+1 == len(source) {
				return nil, false
			}
			next := source[i+1]
			if quote == '\'' && next != '\\' && next != '\'' {
				word.WriteByte(c)
				started = true
				continue
			}
			i++
			switch next {
			case 'c':
				if quote != 0 {
					return nil, false
				}
				flush()
				return words, true
			case '_':
				if quote == 0 {
					flush()
					continue
				}
				c = ' '
			case 'n':
				c = '\n'
			case 'r':
				c = '\r'
			case 't':
				c = '\t'
			case 'v':
				c = '\v'
			case 'f':
				c = '\f'
			case '\\', '\'', '"', '#', '$':
				c = next
			default:
				return nil, false
			}
			word.WriteByte(c)
			started = true
			continue
		}
		if quote != 0 {
			if c == quote {
				quote = 0
			} else {
				word.WriteByte(c)
			}
			started = true
			continue
		}
		switch c {
		case '\'', '"':
			quote = c
			started = true
		case ' ', '\t', '\n', '\r', '\v', '\f':
			flush()
		case '#':
			if !started {
				return words, true
			}
			word.WriteByte(c)
		default:
			word.WriteByte(c)
			started = true
		}
	}
	if quote != 0 {
		return nil, false
	}
	flush()
	return words, true
}

func envInterpreter(source string) string {
	args, ok := splitEnvShebang(source)
	if !ok {
		return ""
	}
	options := true
	for len(args) > 0 {
		arg := args[0]
		if options && arg == "--" {
			options = false
			args = args[1:]
		} else if options && strings.HasPrefix(arg, "--split-string=") {
			args[0] = strings.TrimPrefix(arg, "--split-string=")
			if args[0] == "" {
				args = args[1:]
			}
		} else if options && strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") && strings.Contains(arg, "S") && strings.TrimLeft(strings.SplitN(arg[1:], "S", 2)[0], "i0v") == "" {
			// -S and grouped forms such as -vSpython3 carry the command inline.
			args[0] = strings.SplitN(arg[1:], "S", 2)[1]
			if args[0] == "" {
				args = args[1:]
			}
		} else if options && envOperandArgs.MatchString(arg) {
			if len(args) < 2 {
				return ""
			}
			args = args[2:]
		} else if envVarArgs.MatchString(arg) || (options && envOptArgs.MatchString(arg)) {
			args = args[1:]
		} else {
			return arg
		}
	}
	return ""
}
