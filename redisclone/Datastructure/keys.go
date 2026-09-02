package Datastructure

// while we do have commands like the del removes the or more keys from the store
// the store will handle the background go routines.
// Note: we are using the Lock and s.Unlock to handle concurrency
// We remove the one or more keys from the store
func (s *Store) Del(keys ...string) int {
	// to handle any concurrecny make sure to use Lock and Unlock
	s.mu.Lock()
	defer s.mu.Unlock()
	// we
	deleted := 0
	// the keys string arguments are being collected into a slice
	for _, key := range keys {
		// the getstore check if the store is entry
		if s.getStore(key) != nil {
			delete(s.data, key)
			deleted++ // con
		} // append the key command to the file for data presitence

	}
	return deleted
}

// Get retrieves the value of the key and check if the
func (s *Store) Get(keys ...string) int {
	s.mu.Unlock()
	defer s.mu.Lock() // to handle
	// entry keys to get the key
	getter := 0
	//s.data[keys] will run into an error for the map inde
	for _, ok := range keys {
		if s.getStore(ok) != nil {
			delete(s.data, ok)
			getter++
			s.mu.Lock()
		}
	}
	return getter
}

// Exists return how many of the given keys are present and is they are unexpired
// Note: in redis when it comes to dealing with duplicate keys we count them separately
func (s *Store) Exists(keys ...string) int {
	// again we need to handle concurrency
	s.mu.Lock()
	defer s.mu.Unlock()
	exist := 0
	for _, key := range keys {
		// now we check if the key does exsit
		if s.getStore(key) != nil {
			exist++
		}
	}
	return exist
}

// Type returns the redis type name of the key itself
// None means the key does not exist or has expired
// the type will point the data structures, like the Strings,Lists, Set, Hashes
func (s *Store) Type(key string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	typ := s.getStore(key)
	if typ != nil {
		return typ.typeName() // comes from the type in entry.go file
	} else {
		return "error"
	}
}

// the last two functions check the len of the current keys proccessed in the map/hash
// and check what are the remaining keys left
func (s *Store) remainingKey() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	//	if len(s.data) == 0{
	//		return 0
	//	} if we want to use a stack we can chekc
	return len(s.data)
}

// at the end we do need to delete all the keys
func (s *Store) Flushkeys() {
	s.mu.Lock()
	defer s.mu.Unlock()
	// flush the entire key and value pair in the dict
	s.data = make(map[string]*Entry)
}

// the redis command for dbsize

func (s *Store) DbSize() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.data)
}
