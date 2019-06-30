package ringslice

import (
	"errors"
	"fmt"
)

// Slice struct
type Slice struct {
	values []interface{}
	used   int
	start  int
	end    int
	debug  bool
	cap    int
	wipe   func(int, []interface{})
}

// TODO provide a clear function

// NewSlice does
func NewSlice(capacity int, debug bool, wipe func(int, []interface{})) *Slice {
	return &Slice{values: make([]interface{}, capacity), debug: debug, cap: capacity, wipe: wipe}
}

// Append does
func (s *Slice) Append(value interface{}) error {
	if s.used == s.cap {
		return errors.New("Index is full cannot append")
	}
	ind := s.trueIndex(s.start, s.used) // next index is same as num written
	if s.debug {
		fmt.Println("ind to write to ", ind)
	}
	s.values[ind] = value
	s.used++
	return nil
}

func (s *Slice) Values(value func(interface{}) int64) []int64 {
	v := []int64{}
	for _, val := range s.values {
		v = append(v, value(val))
	}
	return v
}

func (s *Slice) Stats() {
	fmt.Println("used", s.used, "start", s.start, "end", s.end)
}

func (s *Slice) Purge(want int64, value func(interface{}) int64) {
	ind := s.FindClosestBelowOrEqual(want, value)
	if ind == -1 {
		return
	}
	fmt.Println("deleting bounds", s.start, ind)
	s.DeleteBounds(s.start, ind)
}

func (s *Slice) FindClosestBelowOrEqual(want int64, value func(interface{}) int64) int {
	if s.used == 0 {
		return -1
	}
	// binary search and call value on node to find value
	start := s.trueIndex(s.start, 0)
	falseMax := start + s.used - 1 // if the slice were continous, the highest index
	end := falseMax

	if end == start { // there's one node
		if value(s.values[s.trueIndex(start, 0)]) > want {
			return -1 // don't delete it, it's too low
		}
		return start
	}
	count := 0
	for m := (start + falseMax) / 2; ; {
		count++
		fmt.Println(start, end, m)
		if count > 5 {
			panic(count)
		}
		cur := value(s.values[s.trueIndex(m, 0)])
		if cur == want {
			// if we find the exact value, walk to latest index with this value
			return s.findLatestEquivalent(m, want, value)
		}

		// if the start and end indices are 1 away, return higher index <= want
		if end-start == 1 || (end == 0 && start == s.cap-1) {
			return s.determineBoundary(start, end, want, value)
		}

		// if the value we check was less, set the start of next midpoint here
		if cur < want {
			start = m
		}
		// if the value we check was greater, set end to be this point
		if cur > want {
			end = m
		}
		// set the next check index to be midpoint between changed start and end
		m = (start + end) / 2
	}
}

func (s *Slice) findLatestEquivalent(m int, want int64, value func(interface{}) int64) int {
	new := m
	for {
		new = s.next(new)
		if new == m {
			return s.end // we've done a full loop
		}
		if value(s.values[s.trueIndex(new, 0)]) != want {
			return s.prev(new)
		}
	}
}

// give highest index of values below want
func (s *Slice) determineBoundary(start, end int, want int64, value func(interface{}) int64) int {
	if value(s.values[s.trueIndex(end, 0)]) <= want {
		return end
	}
	if value(s.values[s.trueIndex(start, 0)]) > want {
		return -1
	}
	return start

}

// DeleteBounds does
func (s *Slice) DeleteBounds(start, end int) {
	stop := s.trueIndex(end, 0)
	for i := start; ; i = s.next(i) {
		s.wipe(i, s.values)
		if s.used > 0 {
			s.used-- // TODO: could speed this up with calculation
		}
		if i == stop {
			break
		}
	}
	s.start = s.next(stop)
}

// DeleteCount does
func (s *Slice) DeleteCount(count int) {
	ind := s.start
	if count > s.used {
		count = s.used // save us some time
	}
	for i := 0; i < count; i++ {
		s.wipe(i, s.values)
		ind = s.next(ind)
	}
	s.used -= count
	s.start = ind
}

func (s *Slice) next(cur int) int {
	if cur == s.cap-1 {
		return 0
	}
	return cur + 1
}
func (s *Slice) prev(cur int) int {
	if cur == 0 {
		return s.cap - 1
	}
	return cur - 1
}

// returns index of length AWAY from start taking into account wrap around
// start,0 == start
func (s *Slice) trueIndex(start, length int) int {
	return (start + length) % s.cap
}
