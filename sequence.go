package factorybot

// Sequence represents a sequence
type Sequence struct {
	counter int
	value   func(int) interface{}
}

// One returns next value from the sequence
func (s *Sequence) One() interface{} {
	s.counter++

	if s.value != nil {
		return s.value(s.counter)
	}

	return s.N()
}

// N returns next counter value
func (s *Sequence) N() int {
	s.counter++
	return s.counter
}

// Rewind resets the sequence counter
func (s *Sequence) Rewind() {
	s.counter = 0
}

// NewSequence creates new sequence
func NewSequence(params ...func(n int) interface{}) Sequence {
	if len(params) > 0 {
		return Sequence{value: params[0]}
	}

	return Sequence{}
}
