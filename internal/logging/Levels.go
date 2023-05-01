// Package logging implements leveled logging mechanism
package logging

// Level alias to int type.
type Level int

// Levels available logging level.
var Levels = struct {
	ERR   Level
	WARN  Level
	INFO  Level
	DEBUG Level
	OFF   Level
}{
	OFF:   5,
	ERR:   4,
	WARN:  3,
	INFO:  2,
	DEBUG: 1,
}
