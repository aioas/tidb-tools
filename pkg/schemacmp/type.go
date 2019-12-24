package schemacmp

import (
	"github.com/pingcap/parser/mysql"
	"github.com/pingcap/parser/types"
)

// encodeTypeTpAsLattice
func encodeFieldTypeToLattice(ft *types.FieldType) Tuple {
	var flen, dec Lattice
	if ft.Tp == mysql.TypeNewDecimal {
		flen = Singleton(ft.Flen)
		dec = Singleton(ft.Decimal)
	} else {
		flen = Int(ft.Flen)
		dec = Int(ft.Decimal)
	}

	var defVal Lattice
	if mysql.HasAutoIncrementFlag(ft.Flag) || !mysql.HasNoDefaultValueFlag(ft.Flag) {
		defVal = Maybe(Singleton(ft.Flag & (mysql.AutoIncrementFlag | mysql.NoDefaultValueFlag)))
	} else {
		defVal = Maybe(nil)
	}

	return Tuple{
		// TODO: Currently we treat distinct types as incompatible.
		// In reality e.g. TINYINT is compatible with BIGINT.
		Singleton(ft.Tp), // 0
		flen,             // 1
		dec,              // 2

		// TODO: recognize if the remaining flags can be merged or not.
		Singleton(ft.Flag &^ (mysql.AutoIncrementFlag | mysql.NoDefaultValueFlag | mysql.NotNullFlag)), // 3
		Bool(!mysql.HasNotNullFlag(ft.Flag)), // 4 - NON NULL flag (NULL is more compatible than NOT NULL).
		defVal,                               // 5 - AUTO_INCREMENT and DEFAULT flag.

		Singleton(ft.Charset), // 6
		Singleton(ft.Collate), // 7

		StringList(ft.Elems), // 8
	}
}

func decodeFieldTypeFromLattice(tup Tuple) *types.FieldType {
	lst := tup.Unwrap().([]interface{})

	flags := lst[3].(uint)
	if !lst[4].(bool) {
		flags |= mysql.NotNullFlag
	}
	if x, ok := lst[5].(uint); ok {
		flags |= x
	} else {
		flags |= mysql.NoDefaultValueFlag
	}

	return &types.FieldType{
		Tp:      lst[0].(byte),
		Flen:    lst[1].(int),
		Decimal: lst[2].(int),
		Flag:    flags,
		Charset: lst[6].(string),
		Collate: lst[7].(string),
		Elems:   lst[8].([]string),
	}
}

type typ struct{ Tuple }

func Type(ft *types.FieldType) Lattice {
	return typ{Tuple: encodeFieldTypeToLattice(ft)}
}

func (a typ) Unwrap() interface{} {
	return decodeFieldTypeFromLattice(a.Tuple)
}

func (a typ) Compare(other Lattice) (int, error) {
	if b, ok := other.(typ); ok {
		return a.Tuple.Compare(b.Tuple)
	}
	return 0, IncompatibleError{}
}

func (a typ) Join(other Lattice) (Lattice, error) {
	if b, ok := other.(typ); ok {
		return a.Tuple.Join(b.Tuple)
	}
	return nil, IncompatibleError{}
}
