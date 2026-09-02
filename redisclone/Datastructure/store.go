package Datastructure

import (
	"fmt"
	"os"

	"strconv"
	"sync"
)

// note redis is a data store. the Store struct is a in memory data related structures
// RwMutex = multiple client goroutines run simultaneously, so we use Mutex to block many concurrent calls
// Store is the central in memory key-value database
type Store struct {
	mu   sync.RWMutex //
	data map[string]*Entry
	// we append only the log file for data presistence
	aofFile os.File
}

// Now we create a way to initializes the store and start the background time to live
// we initialize the store with which loads in peresisted data -> Redis AOF
func NewStore() *Store {
	s := &Store{
		data: make(map[string]*Entry),
	}
	//load the presitence for the append only file. It goes in this background

	// then we call a goroutine. a single thread has so many go routines
	// we will run the
	go s.activeExpiryLoop()
	return s

}

// get returns the live entry for key, or nil if the key does not exisit or expired
// Note: for myself this will be called for each time we get a key when we del a key
func (s *Store) getStore(key string) *Entry {
	element, ok := s.data[key]
	// check if the key was even set
	if !ok {
		return nil
	}
	// call the element
	if element.isExpired() {
		return nil
	}
	return element
}

// it takes a s as a parameter and
func mustInt(s string) (int, error) {
	val, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Print("Error Occured")
	}
	return val, nil

}
