package util

import "strings"

type CmdLines []string

func RemoveQInLines(lines []string) (res []string) {
	for i, line := range lines {
		if line == "qman" {
			line = "man"
		}
		if line == "qtop" {
			line = "top"
		}
		lines[i] = line
	}

	res = []string{}
	for i, line := range lines {
		if strings.HasPrefix(line, "q") && i-1 >= 0 {
			preCmd := lines[i-1]
			if preCmd == "top" || preCmd == "man" || strings.HasPrefix(preCmd, "man ") {
				line = strings.TrimPrefix(line, "q")
			}
		}
		res = append(res, line)
	}
	return res
}
