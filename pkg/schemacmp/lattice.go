package schemacmp

type IncompatibleError struct{}

func (IncompatibleError) Error() string {
	return "the inputs are incompatible"
}

// Lattice is implemented for types which forms a join-semilattice.
type Lattice interface {
	// Unwrap returns the underlying object supporting the lattice. This
	// operation is deep.
	Unwrap() interface{}

	// Compare this instance with another instance.
	//
	// Returns -1 if `self < other`, 0 if `self == other`, 1 if `self > other`.
	// Returns `IncompatibleError` if the two instances are not ordered.
	Compare(other Lattice) (int, error)

	// Join finds the "least upper bound" of two Lattice instances. The result
	// is `>=` both inputs. Returns an error if the join does not exist.
	Join(other Lattice) (Lattice, error)
}

// Bool is a boolean implementing Lattice where `false < true`.
type Bool bool

// Unwrap implements Lattice
func (a Bool) Unwrap() interface{} {
	return bool(a)
}

// Compare implements Lattice.
func (a Bool) Compare(other Lattice) (int, error) {
	b, ok := other.(Bool)
	switch {
	case !ok:
		return 0, IncompatibleError{}
	case a == b:
		return 0, nil
	case bool(a):
		return 1, nil
	default:
		return -1, nil
	}
}

// Join implements Lattice
func (a Bool) Join(other Lattice) (Lattice, error) {
	b, ok := other.(Bool)
	if !ok {
		return nil, IncompatibleError{}
	}
	return a || b, nil
}

type singleton struct{ value interface{} }

// Unwrap implements Lattice
func (a singleton) Unwrap() interface{} {
	return a.value
}

// Singleton wraps an unordered value. Distinct instances of Singleton are
// incompatible.
func Singleton(value interface{}) Lattice {
	return singleton{value: value}
}

// Compare implements Lattice.
func (a singleton) Compare(other Lattice) (int, error) {
	if b, ok := other.(singleton); ok && a.value == b.value {
		return 0, nil
	}
	return 0, IncompatibleError{}
}

// Join implements Lattice
func (a singleton) Join(other Lattice) (Lattice, error) {
	if b, ok := other.(singleton); ok && a.value == b.value {
		return a, nil
	}
	return nil, IncompatibleError{}
}

// BitSet is a set of bits where `a < b` iff `a` is a subset of `b`.
type BitSet uint

// Unwrap implements Lattice.
func (a BitSet) Unwrap() interface{} {
	return uint(a)
}

// Compare implements Lattice.
func (a BitSet) Compare(other Lattice) (int, error) {
	b, ok := other.(BitSet)
	switch {
	case !ok:
		return 0, IncompatibleError{}
	case a == b:
		return 0, nil
	case a&^b == 0:
		return -1, nil
	case b&^a == 0:
		return 1, nil
	default:
		return 0, IncompatibleError{}
	}
}

// Join implements Lattice.
func (a BitSet) Join(other Lattice) (Lattice, error) {
	b, ok := other.(BitSet)
	if !ok {
		return nil, IncompatibleError{}
	}
	return a | b, nil
}

// Byte is a byte implementing Lattice.
type Byte byte

// Unwrap implements Lattice.
func (a Byte) Unwrap() interface{} {
	return byte(a)
}

// Compare implements Lattice.
func (a Byte) Compare(other Lattice) (int, error) {
	b, ok := other.(Byte)
	switch {
	case !ok:
		return 0, IncompatibleError{}
	case a == b:
		return 0, nil
	case a > b:
		return 1, nil
	default:
		return -1, nil
	}
}

// Join implements Lattice.
func (a Byte) Join(other Lattice) (Lattice, error) {
	b, ok := other.(Byte)
	switch {
	case !ok:
		return nil, IncompatibleError{}
	case a >= b:
		return a, nil
	default:
		return b, nil
	}
}

// Int is an int implementing Lattice.
type Int int

// Unwrap implements Lattice.
func (a Int) Unwrap() interface{} {
	return int(a)
}

// Compare implements Lattice.
func (a Int) Compare(other Lattice) (int, error) {
	b, ok := other.(Int)
	switch {
	case !ok:
		return 0, IncompatibleError{}
	case a == b:
		return 0, nil
	case a > b:
		return 1, nil
	default:
		return -1, nil
	}
}

// Join implements Lattice.
func (a Int) Join(other Lattice) (Lattice, error) {
	b, ok := other.(Int)
	switch {
	case !ok:
		return nil, IncompatibleError{}
	case a >= b:
		return a, nil
	default:
		return b, nil
	}
}

// Meet implements Lattice.
func (a Int) Meet(other Lattice) (Lattice, error) {
	b, ok := other.(Int)
	switch {
	case !ok:
		return nil, IncompatibleError{}
	case a <= b:
		return a, nil
	default:
		return b, nil
	}
}

// Tuple of Lattice instances. Given two Tuples `a` and `b`, we define `a < b`
// iff `a[i] < b[i]` for all `i`.
type Tuple []Lattice

// Unwrap implements Lattice. The returned type is a `[]interface{}`.
func (a Tuple) Unwrap() interface{} {
	res := make([]interface{}, 0, len(a))
	for _, value := range a {
		res = append(res, value.Unwrap())
	}
	return res
}

// Compare implements Lattice.
func (a Tuple) Compare(other Lattice) (int, error) {
	b, ok := other.(Tuple)
	if !ok || len(b) != len(a) {
		return 0, IncompatibleError{}
	}
	result := 0
	for i, left := range a {
		res, err := left.Compare(b[i])
		if err != nil {
			return 0, err
		}
		result, err = CombineCompareResult(result, res)
		if err != nil {
			return 0, err
		}
	}
	return result, nil
}

// CombineCompareResult combines two comparison results.
func CombineCompareResult(x int, y int) (int, error) {
	switch {
	case x == y || y == 0:
		return x, nil
	case x == 0:
		return y, nil
	default:
		return 0, IncompatibleError{}
	}
}

// Join implements Lattice
func (a Tuple) Join(other Lattice) (Lattice, error) {
	b, ok := other.(Tuple)
	if !ok || len(b) != len(a) {
		return nil, IncompatibleError{}
	}
	result := make(Tuple, 0, len(a))
	for i, left := range a {
		res, err := left.Join(b[i])
		if err != nil {
			return nil, err
		}
		result = append(result, res)
	}
	return result, nil
}

type maybe struct{ Lattice }

// Maybe includes `nil` as the universal lower bound of the original Lattice.
func Maybe(inner Lattice) Lattice {
	return maybe{Lattice: inner}
}

// Unwrap implements Lattice.
func (a maybe) Unwrap() interface{} {
	if a.Lattice != nil {
		return a.Lattice.Unwrap()
	}
	return nil
}

// Compare implements Lattice.
func (a maybe) Compare(other Lattice) (int, error) {
	b, ok := other.(maybe)
	switch {
	case !ok:
		return 0, IncompatibleError{}
	case a.Lattice == nil && b.Lattice == nil:
		return 0, nil
	case a.Lattice == nil:
		return -1, nil
	case b.Lattice == nil:
		return 1, nil
	default:
		return a.Lattice.Compare(b.Lattice)
	}
}

// Join implements Lattice.
func (a maybe) Join(other Lattice) (Lattice, error) {
	b, ok := other.(maybe)
	switch {
	case !ok:
		return nil, IncompatibleError{}
	case a.Lattice == nil:
		return b, nil
	case b.Lattice == nil:
		return a, nil
	default:
		join, err := a.Lattice.Join(b.Lattice)
		if err != nil {
			return nil, err
		}
		return maybe{Lattice: join}, nil
	}
}

// StringList is a list of string where `a <= b` iff `a == b[:len(a)]`.
type StringList []string

// Unwrap implements Lattice.
func (a StringList) Unwrap() interface{} {
	return []string(a)
}

// Compare implements Lattice.
func (a StringList) Compare(other Lattice) (int, error) {
	b, ok := other.(StringList)
	if !ok {
		return 0, IncompatibleError{}
	}
	minLen := len(a)
	if minLen > len(b) {
		minLen = len(b)
	}
	for i := 0; i < minLen; i++ {
		if a[i] != b[i] {
			return 0, IncompatibleError{}
		}
	}
	switch {
	case len(a) == len(b):
		return 0, nil
	case len(a) < len(b):
		return -1, nil
	default:
		return 1, nil
	}
}

// Join implements Lattice.
func (a StringList) Join(other Lattice) (Lattice, error) {
	cmp, err := a.Compare(other)
	switch {
	case err != nil:
		return nil, err
	case cmp <= 0:
		return other, nil
	default:
		return a, nil
	}
}

// LatticeMap is a map of Lattice objects keyed by strings.
type LatticeMap interface {
	// Unwrap returns the underlying object supporting the LatticeMap.
	Unwrap() interface{}

	// New creates an empty LatticeMap of the same type as the receiver.
	New() LatticeMap

	// Insert inserts a key-value pair into the map.
	Insert(key string, value Lattice)

	// Get obtains the Lattice object at the given key. Returns nil if the key
	// does not exist.
	Get(key string) Lattice

	// ForEach iterates the map.
	ForEach(func(key string, value Lattice) error) error

	// CompareWithNil returns the comparison result when the value is compared
	// with a non-existing entry.
	CompareWithNil(value Lattice) (int, error)

	// JoinWithNil returns the result when the value is joined with a
	// non-existing entry. If the joined result should be non-existing, this
	// method should return nil, nil.
	JoinWithNil(value Lattice) (Lattice, error)

	// ShouldDeleteIncompatibleJoin returns true if two incompatible entries
	// should be deleted instead of propagating the error.
	ShouldDeleteIncompatibleJoin() bool
}

type latticeMap struct{ LatticeMap }

// Unwrap implements Lattice.
func (a latticeMap) Unwrap() interface{} {
	return a.LatticeMap.Unwrap()
}

func (a latticeMap) iter(other Lattice, action func(k string, av, bv Lattice) error) error {
	b, ok := other.(latticeMap)
	if !ok {
		return IncompatibleError{}
	}

	visitedKeys := make(map[string]struct{})
	err := a.ForEach(func(k string, av Lattice) error {
		visitedKeys[k] = struct{}{}
		return action(k, av, b.Get(k))
	})
	if err != nil {
		return err
	}

	return b.ForEach(func(k string, bv Lattice) error {
		if _, ok := visitedKeys[k]; ok {
			return nil
		}
		return action(k, nil, bv)
	})
}

// Compare implements Lattice.
func (a latticeMap) Compare(other Lattice) (int, error) {
	result := 0
	err := a.iter(other, func(k string, av, bv Lattice) error {
		var (
			cmpRes int
			e      error
		)
		switch {
		case av != nil && bv != nil:
			cmpRes, e = av.Compare(bv)
		case av != nil:
			cmpRes, e = a.CompareWithNil(av)
		default:
			cmpRes, e = a.CompareWithNil(bv)
			cmpRes = -cmpRes
		}
		if e != nil {
			return e
		}
		result, e = CombineCompareResult(result, cmpRes)
		return e
	})
	return result, err
}

// Join implements Lattice.
func (a latticeMap) Join(other Lattice) (Lattice, error) {
	result := a.New()
	err := a.iter(other, func(k string, av, bv Lattice) error {
		var (
			joinRes Lattice
			e       error
		)
		switch {
		case av != nil && bv != nil:
			joinRes, e = av.Join(bv)
		case av != nil:
			joinRes, e = a.JoinWithNil(av)
		default:
			joinRes, e = a.JoinWithNil(bv)
		}
		if e != nil {
			if a.ShouldDeleteIncompatibleJoin() {
				return nil
			}
			return e
		}
		result.Insert(k, joinRes)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return latticeMap{LatticeMap: result}, nil
}

// Map wraps a LatticeMap instance into a Lattice.
func Map(lm LatticeMap) latticeMap {
	return latticeMap{LatticeMap: lm}
}
