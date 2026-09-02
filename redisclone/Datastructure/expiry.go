package Datastructure

import (
	"time"
)

//the expiry go will handles any background go routine that sweeps up any reamining keys
//For future reference, this is perfect for redis acting as cache

// the go routine
// Redis uses  Passive Expiration to handle clients access a key if the key Time To Live
// this happens becuse redis needs to know when to expend keys and when to hold them
// run forver in the background
func (s *Store) activeExpiryLoop() {
	ticker := time.NewTicker(100 * time.Millisecond)
	// defer closes the connection after the file has ran
	defer ticker.Stop()
	for range ticker.C {
		// for concurrency issues: redis uses one thread but stops it from the
		s.mu.Lock()
		for key, e := range s.data {
			if e.isExpired() {
				delete(s.data, key)
			}
		}
		s.mu.Unlock()
	}
}
