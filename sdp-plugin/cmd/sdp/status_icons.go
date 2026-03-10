package main

func boolIcon(value bool) string {
	if value {
		return "[OK]"
	}
	return "[MISSING]"
}
