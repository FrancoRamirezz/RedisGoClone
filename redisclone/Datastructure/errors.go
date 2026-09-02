package Datastructure

import "errors"

// You can just write the errors in their own file, but this makes it easier to implement
// here we will check for three instance : one where the command is used on the wrong key
var ErrWrongType = errors.New("Wrong type operation: The key is holding the wrong kind of value")

var ErrNotIntger = errors.New("Err value not the correct integer, or out of inputs range")

var ErrKeyNotFound = errors.New("Err key occured")
