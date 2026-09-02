package Datastructure

// a lists in redis are bit different from python
// they follow more a linked list where we can
// we can push to the list, starting from the front

// Model a que using RPush and LPop
// If we want to retrieve elements with LRange
// Get the length of a list with LLEN
func (s *Store) LPush(keys string, values ...string) (int, error) {
	// as before since we are using goroutines we can deal with concurrecny
	s.mu.Lock()
	defer s.mu.Unlock()
	// this means were about to push to the list
	push := s.getStore(keys)
	// check if were about to push to the left meaning its first one to be added
	// we set the array lists to be added
	// s.data[key] = Entry{Value: value,} follow this data formula
	if push == nil {
		push = &Entry{typ: typeList, list: []string{}}
		s.data[keys] = push
	} else if push.typ != typeList {
		return 0, ErrWrongType // comes from the helper function side
	}
	// append the values into the list
	// note when we use for index, value := range
	for _, value := range values {
		push.list = append([]string{value}, push.list...)
		// Note the
	}
	return len(push.list), nil

}

// now we push on to the right side
func (s *Store) RPush(keys string, values ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	push := s.getStore(keys)
	if push == nil {
		push = &Entry{typ: typeList, list: []string{}}
		s.data[keys] = push
	} else if push.typ != typeList {
		return 0, ErrWrongType
	}
	push.list = append(push.list, values...)
	return len(push.list), nil
}

// Lpop follows the que method where we will pop from the left
func (s *Store) LPop(key string) (string, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	pop := s.getStore(key)
	if pop == nil || len(pop.list) == 0 {
		return "", false, nil
	}
	if pop.typ != typeList {
		return "", false, ErrWrongType
	}
	ind := pop.list[0]      //index for the front
	pop.list = pop.list[1:] // pop from the front following the que
	return ind, true, nil
}

func (s *Store) RPop(key string) (string, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	pop := s.getStore(key)
	if pop == nil || len(pop.list) == 0 {
		return "", false, ErrWrongType
	}
	if pop.typ != typeList {
		return "", false, ErrWrongType
	}
	n := len(pop.list)
	indx := pop.list[n-1]
	pop.list = pop.list[:n-1]
	return indx, true, nil
}

// the len returns the number of elements in the list
// return zero if the element does not exisit
func (s *Store) LLen(key string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	elem := s.getStore(key)
	if elem == nil || elem.typ != typeList {
		return 0, ErrWrongType
	}
	return len(elem.list), nil
}

// LRange returns a slice of the list between start and stop
// Supports negative indices -1 last element
func (s *Store) LRange(key string, start, stop int) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	elem := s.getStore(key)
	if elem == nil || elem.typ != typeList {
		return nil, ErrWrongType
		// note i am not making
	}
	n := len(elem.list)
	if n == 0 {
		return []string{}, nil
	} // convert negative indices to positive
	if start < 0 {
		start = n + start
	}
	if stop < 0 {
		stop = n + stop
	}
	// conert negative indices to positive
	// Clamp to valid range
	if start < 0 {
		start = 0
	}
	// start with a empty range
	if start > stop {
		return []string{}, nil
	}
	// make a var
	res := make([]string, stop-start+1)
	return res, nil
}
