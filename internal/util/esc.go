package util

// EscReturns changes any actual carriage returns or line returns into their
// backslashed equivalents.
func EscReturns(s string) string {
	var out string
	for _, r := range s {
		switch r {
		case '\n':
			out += "\\n"
		case '\r':
			out += "\\r"
		default:
			out += string(r)
		}
	}
	return out
}

// UnEscReturns changes any escaped carriage returns or line returns into their
// actual values.
func UnEscReturns(s string) string {
	var out string
	n := len(s)
	for i := 0; i < n; i++ {
		if s[i] == '\\' && i != n-1 {
			if s[i+1] == 'n' {
				out += "\n"
				i++
				continue
			}
			if s[i+1] == 'r' {
				out += "\r"
				i++
				continue
			}
		}
		out += string(s[i])
	}
	return out
}
