package logging

// Color list of predefined color log logging level.
var Color = struct {
	RED    string
	ORANGE string
	GREEN  string
	BLUE   string
	CYAN   string
	CLEAR  string
}{
	RED:    "\x1b[0;31m",
	ORANGE: "\x1b[0;33m",
	GREEN:  "\x1b[0;32m",
	BLUE:   "\x1b[0;34m",
	CYAN:   "\x1b[0;36m",
	CLEAR:  "\x1b[0m",
}
