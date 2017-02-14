package logger

import "sync"
import "strings"

// FilterByMinimumLevel only allows messages that are equal to or higher than the specified minimum level
func FilterByMinimumLevel(levelFilter *LevelFilter) FilterFunc {

	return levelFilter.filter
}

// LogLevel is a string defining a specific level of log messages
type LogLevel string

// LevelFilter is used to filter out messages that don't meet the specified minimum level
//
// Once the filter is in use somewhere, it is not safe to modify the structure.
type LevelFilter struct {
	// Levels is the list of log levels, in increasing order of
	// severity. Example might be: {"Debug", "Info", "Warning", "Error"}.
	Levels []LogLevel

	// MinLevel is the minimum level allowed through
	MinLevel LogLevel

	badLevels map[LogLevel]struct{}
	once      sync.Once
}

// NewLevelFilter creates a new filter based on the given minimum level
//
// Note that the available levels do not need to be specified, and are set to
// what this library is capable of.
func NewLevelFilter(minLevel LogLevel) *LevelFilter {
	return &LevelFilter{
		Levels:   []LogLevel{"Debug", "Info", "Warning", "Error"},
		MinLevel: minLevel,
	}
}

// SetMinLevel is used to update the minimum log level
func (f *LevelFilter) SetMinLevel(min LogLevel) {
	f.MinLevel = min
	f.init()
}

// init calculates the disallowed levels, to prepare for easier filtering
func (f *LevelFilter) init() {
	badLevels := make(map[LogLevel]struct{})
	for _, level := range f.Levels {
		if strings.EqualFold(string(level), string(f.MinLevel)) {
			break
		}
		badLevels[level] = struct{}{}
	}
	f.badLevels = badLevels
}

func (f *LevelFilter) filter(e *Entry) bool {
	f.once.Do(f.init)

	// No entry? Then allow.
	if e == nil {
		return true
	}

	// Check if it exists in the list of bad level
	level := LogLevel(e.Level)

	_, ok := f.badLevels[level]
	return !ok
}
