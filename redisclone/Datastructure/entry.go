package Datastructure

import "time"

// Overview: this file handles the time to live key for our commands that communicate through the resp
// entryType identifies which data type a key holds
// the map in the Store is map[string]*entry
type entryType int

const (
	typeString entryType = iota
	typeHash             // field value map
	typeList             // ordered list of strings
	typeSet
)

// handles the internal representation of a stored value
// the map is in map[string]*Entry
type Entry struct {
	// this handles more than a single type
	// our handles str hash, a list, set,
	typ    entryType
	str    string
	hash   map[string]string
	list   []string
	set    map[string]struct{}
	expiry time.Time // gives an error once the
}

// the function isExpired returns True
// Note: this function will be called
func (e *Entry) isExpired() bool {
	// the zero time means no expiry was set
	// the time to live is not set
	if e.expiry == 0 {
		return false
	}
	return time.Now().After(e.expiry)
}

// for our redis we want to develop different data structures for it
// Used by the Type Command
func (e *Entry) typeName() string {
	switch e.typ {
	case typeString:
		return "string"
	case typeHash:
		return "hash"
	case typeList:
		return "list"
	case typeSet:
		return "set"
	default:
		return "None"
	}
}
