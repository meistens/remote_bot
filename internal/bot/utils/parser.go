package utils

import "strings"

// parse command args
func ParseCommand(text string) (command string, args map[string]string) {
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return "", nil
	}

	command = parts[0]
	args = make(map[string]string)

	for i := 1; i < len(parts); i++ {
		if strings.HasPrefix(parts[i], "--") {
			key := strings.TrimPrefix(parts[i], "--")
			if i+1 < len(parts) && !strings.HasPrefix(parts[i+1], "--") {
				args[key] = parts[i+1]
				i++ // Skip the value
			} else {
				args[key] = "true"
			}
		}
	}
	return command, args
}
