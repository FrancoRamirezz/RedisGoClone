package Datastructure

// Redis uses a hash as a data structure that stores a set of field value pairs under a single key
// set one or more fields in the hash at the key
func (s *Store) HSet(key string, pairs map[string]string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock() // lock the store shared between multiple goroutines
	set := s.getStore(key)
	// we check if the key does not exist, if so we will create new hash entry
	if set == nil {
		set = &Entry{
			typ:  typeHash,
			hash: make(map[string]string), // we create the actual hash that will hold fields and values
		}
		s.data[key] = set
	}
	added := 0 // in redis we need to keep track of the number of new fields that are addeded
	// make sure to loop through the fields
	for k, v := range pairs {
		// check if the whether the field already exists using maps in go
		if _, exists := set.hash[k]; !exists {
			added++
		}
		set.hash[k] = v
	}
	return added, nil // at the end make sure to return how many fields we created
}

// make a get function that returns a single field from the hash at the key
func (s *Store) HGet(key string, field string) (string, bool, error) {
	// as above we need to handle any concurrenct
	s.mu.Lock()
	defer s.mu.Unlock()
	get := s.getStore(key)
	if get == nil {
		// means we have a missing key
		return "", false, nil
	}
	if get.typ != typeHash {
		return "", false, ErrWrongType
	}
	// before we had a key and value pair, but here we will check if the value itself exsits
	val, exists := get.hash[field]
	return val, exists, nil
}

// the hash get all will return all the field and value pairs in the hash
// edge case: if the map does not exsit then we return an empty map
func (s *Store) HGetall(key string) (map[string]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	Hget := s.getStore(key)
	if Hget == nil {
		return map[string]string{}, nil
	}
	if Hget.typ != typeHash {
		return nil, ErrWrongType
	}
	// return the copy of callers
	// we create a new map and change the
	res := make(map[string]string, len(Hget.hash))
	for k, v := range Hget.hash {
		res[k] = v
	}
	return res, nil
}

// HDel deletes one or more fields from the hash
func (s *Store) HDel(key string, fields ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	Hdel := s.getStore(key)
	if Hdel == nil && Hdel.typ != typeHash {
		return 0, ErrWrongType
	}
	deleted := 0
	for _, field := range fields {
		if _, exists := Hdel.hash[field]; exists {
			delete(Hdel.hash, field)
			deleted++
		}
	}
	return deleted, nil
}

// check if the field exists
func (s *Store) HExists(key string, field string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	exist := s.getStore(key)
	// Note
	if exist == nil && exist.typ != typeHash {
		return false, ErrWrongType
	}
	_, exists := exist.hash[field]
	return exists, nil
}

func (s *Store) HLen(key string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	elem := s.getStore(key)
	if elem == nil {
		return 0, nil
	}
	if elem.typ != typeHash {
		return 0, ErrWrongType
	}
	return len(elem.hash), nil
}
