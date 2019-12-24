package schemacmp_test

import (
	. "github.com/pingcap/tidb-tools/pkg/schemacmp"

	. "github.com/pingcap/check"
)

type compatibilitySchema struct{}

var _ = Suite(&compatibilitySchema{})

func (*compatibilitySchema) TestIncompatibleErrorString(c *C) {
	c.Assert(IncompatibleError{}, ErrorMatches, "the inputs are incompatible")
}

func (*compatibilitySchema) TestCompatibilities(c *C) {
	testCases := []struct {
		a             Lattice
		b             Lattice
		compareResult int
		compareError  error
		join          Lattice
		joinError     error
	}{
		{
			a:             Bool(false),
			b:             Bool(false),
			compareResult: 0,
			join:          Bool(false),
		},
		{
			a:             Bool(false),
			b:             Bool(true),
			compareResult: -1,
			join:          Bool(true),
		},
		{
			a:             Bool(true),
			b:             Bool(true),
			compareResult: 0,
			join:          Bool(true),
		},
		{
			a:             Singleton(123),
			b:             Singleton(123),
			compareResult: 0,
			join:          Singleton(123),
		},
		{
			a:            Singleton(123),
			b:            Singleton(2468),
			compareError: IncompatibleError{},
			joinError:    IncompatibleError{},
		},
		{
			a:            BitSet(0b010110),
			b:            BitSet(0b110001),
			compareError: IncompatibleError{},
			join:         BitSet(0b110111),
		},
		{
			a:             BitSet(0xffffffff),
			b:             BitSet(0),
			compareResult: 1,
			join:          BitSet(0xffffffff),
		},
		{
			a:             BitSet(0b10001),
			b:             BitSet(0b11011),
			compareResult: -1,
			join:          BitSet(0b11011),
		},
		{
			a:             BitSet(0x522),
			b:             BitSet(0x522),
			compareResult: 0,
			join:          BitSet(0x522),
		},
		{
			a:             Byte(123),
			b:             Byte(123),
			compareResult: 0,
			join:          Byte(123),
		},
		{
			a:             Byte(1),
			b:             Byte(23),
			compareResult: -1,
			join:          Byte(23),
		},
		{
			a:             Byte(123),
			b:             Byte(45),
			compareResult: 1,
			join:          Byte(123),
		},
		{
			a:            Tuple{Byte(123), Bool(false)},
			b:            Tuple{Byte(67), Bool(true)},
			compareError: IncompatibleError{},
			join:         Tuple{Byte(123), Bool(true)},
		},
		{
			a:             Tuple{},
			b:             Tuple{},
			compareResult: 0,
			join:          Tuple{},
		},
		{
			a:            Tuple{Singleton(6), Singleton(7)},
			b:            Tuple{Singleton(6), Singleton(8)},
			compareError: IncompatibleError{},
			joinError:    IncompatibleError{},
		},
		{
			a:            Tuple{},
			b:            Tuple{Bool(false)},
			compareError: IncompatibleError{},
			joinError:    IncompatibleError{},
		},
		{
			a:            Bool(false),
			b:            Singleton(false),
			compareError: IncompatibleError{},
			joinError:    IncompatibleError{},
		},
		{
			a:            Maybe(Singleton(123)),
			b:            Maybe(Singleton(678)),
			compareError: IncompatibleError{},
			joinError:    IncompatibleError{},
		},
		{
			a:             Maybe(Byte(111)),
			b:             Maybe(Byte(222)),
			compareResult: -1,
			join:          Maybe(Byte(222)),
		},
		{
			a:             Maybe(nil),
			b:             Maybe(Singleton(135)),
			compareResult: -1,
			join:          Maybe(Singleton(135)),
		},
		{
			a:             Maybe(nil),
			b:             Maybe(nil),
			compareResult: 0,
			join:          Maybe(nil),
		},
		{
			a:            Bool(false),
			b:            Maybe(Bool(false)),
			compareError: IncompatibleError{},
			joinError:    IncompatibleError{},
		},
		{
			a:             StringList{"one", "two", "three"},
			b:             StringList{"one", "two", "three", "four", "five"},
			compareResult: -1,
			join:          StringList{"one", "two", "three", "four", "five"},
		},
		{
			a:            StringList{"a", "b", "c"},
			b:            StringList{"a", "e", "i", "o", "u"},
			compareError: IncompatibleError{},
			joinError:    IncompatibleError{},
		},
		{
			a:             StringList{},
			b:             StringList{},
			compareResult: 0,
			join:          StringList{},
		},
	}

	for _, tc := range testCases {
		assert := func(obtained interface{}, checker Checker, args ...interface{}) {
			args = append(args, Commentf("test case = %+v", tc))
			c.Assert(obtained, checker, args...)
		}

		cmp, err := tc.a.Compare(tc.b)
		if tc.compareError != nil {
			assert(err, Equals, tc.compareError)
		} else {
			assert(err, IsNil)
			assert(cmp, Equals, tc.compareResult)
		}

		cmp, err = tc.b.Compare(tc.a)
		if tc.compareError != nil {
			assert(err, Equals, tc.compareError)
		} else {
			assert(err, IsNil)
			assert(cmp, Equals, -tc.compareResult)
		}

		join, err := tc.a.Join(tc.b)
		if tc.joinError != nil {
			assert(err, Equals, tc.joinError)
		} else {
			assert(err, IsNil)
			assert(tc.join, DeepEquals, join)
		}

		join, err = tc.b.Join(tc.a)
		if tc.joinError != nil {
			assert(err, Equals, tc.joinError)
		} else {
			assert(err, IsNil)
			assert(tc.join, DeepEquals, join)
		}

		if tc.joinError == nil {
			cmp, err = join.Compare(tc.a)
			assert(err, IsNil)
			assert(cmp, GreaterEqual, 0)

			cmp, err = join.Compare(tc.b)
			assert(err, IsNil)
			assert(cmp, GreaterEqual, 0)
		}
	}
}
