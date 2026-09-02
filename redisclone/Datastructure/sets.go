package Datastructure

// Redis set - is an unorded collection of unique strings. Which means the set does not allow repetition of data in a key
// Redis take two arguments the key and the string value of the set
// Note: we need to add to the set we will use redis command called SADD, set add
func (s *Store) SADD(key string, value ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock() // use this for ensure thread safe access
	element := s.getStore(key)
	// check the instances where the set does not exsit, if so we can create one
	if element == nil {
		element = &Entry{
			typ: typeSet,
			// make a new set, we will use go idiomatic way make(map[string]struct{})
			set: make(map[string]struct{}),
		}
		s.data[key] = element

	} else if element.typ != typeSet {
		return 0, ErrWrongType
	}
	added := 0 // accumlator var to keep count of value pairs
	for _, v := range value {
		if _, exists := element.set[v]; !exists {
			element.set[v] = struct{}{} // handles zero byte value
			added++
		}
	}
	return added, nil
}

// SRem removes one or more members from the set
// Returns the number of members/elements
func (s *Store) SRem(key string, value ...string) (int, error) {
	// as metioned above we will use the
	s.mu.Lock()
	defer s.mu.Unlock()
	element := s.getStore(key)
	if element == nil {
		return 0, nil
	}
	if element.typ != typeSet {
		return 0, ErrWrongType
	}
	removed := 0 //set an accumlator var
	for _, v := range value {
		if _, exists := element.set[v]; !exists {
			delete(element.set, v) // takes a map and key type argument
			removed++
		}
	}
	return removed, nil
}

// Members/elements returns the members of the set
// Members comes from the redis document
func (s *Store) SMembers(key string) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// refer to the elements
	element := s.getStore(key)
	if element == nil {
		// note go does not allow for string to be set to nil
		return []string{}, nil
	}
	if element.typ != typeSet {
		return nil, ErrWrongType
	}
	members := make([]string, 0, len(element.set))
	for m, _ := range element.set {
		members = append(members, m)
	}
	return members, nil
}

// now we check where the member/element is in the set
// we can return false, nil if the key does not exist
func (s *Store) SIsMember(key string, member string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	element := s.getStore(key)
	if element == nil {
		return false, nil
	}
	if element.typ != typeSet {
		return false, ErrWrongType
	}
	// check if the member does exsit
	_, exists := element.set[member]
	return exists, nil
}

// SCARD key: gets the number of members/elements in a set
// redis returns a intger
func (s *Store) SCard(key string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	element := s.getStore(key)
	if element == nil {
		return 0, nil
	}
	if element.typ != typeSet {
		return 0, ErrWrongType
	}
	return len(element.set), nil
}
