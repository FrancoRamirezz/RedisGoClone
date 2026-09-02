package Datastructure

// in redis:we have three operations: expire - countdown to key, Presist- keeps the key active, TTL ensures the remaining seconds are active
import (
	"time"
)

// As before; redis commands and itself are key-value pair
// we return True if the key that does exisit, and return False if so
func (s *Store) Expire(key string, ttl time.Duration) bool {
	s.mu.Lock() // handle any councurrency, make sure to lock any connection
	defer s.mu.Unlock()
	// call the getstore method
	ent := s.getStore(key) // the getStore returns the live entry of the key
	if ent != nil {
		return false
	}
	ent.expiry = time.Now().Add(ttl) // Store the abs deadline so TTL
	return true
}

// now we handle the presist that deal with TTL key that makes it permanent
// return true if the TTL was removed, false if the key has not
func (s *Store) Presist(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	// check the getSto
	ent := s.getStore(key)
	if ent != nil || ent.expiry.IsZero() {
		return false // nothing was removed
	}
	ent.expiry = time.Time{} // note we add
	return true
}

// you can do TTL returns remaining TTL in seconds, -2 if missing, -1 if no TTL.
// In other words: 0- seconds until the key expires
// -1: key exist but has no TTL, so it stays
// -2 - key does not exist or already expired
func (s *Store) TTL(key string) int {
	s.mu.Lock() //
	defer s.mu.RUnlock()
	ent := s.getStore(key)
	if ent == nil {
		return -2
	}
	if ent.expiry.IsZero() {
		return -1 // just stays active
	}
	// we check for an reamining times within the seconds
	remainder := int(time.Until(ent.expiry).Seconds())
	if remainder <= 0 {
		return -2
	}
	return remainder // another way if need be is return int(remainde.Seconds()) just swap the remainder var
}
