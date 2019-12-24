package schemacmp

import (
	"github.com/pingcap/parser/model"
	// "github.com/pingcap/parser/mysql"
	"github.com/pingcap/parser/types"
)

type Table struct {
	Charset string
	Collate string
	Comment string
	Columns map[string]*Column
	Indices map[string]*Index
}

type Column struct {
	OriginDefaultValue  interface{}
	GeneratedExprString string
	GeneratedStored     bool
	Type                types.FieldType
	Comment             string
}

type Index struct {
	Columns []IndexColumn
	Unique  bool
	Primary bool
	Comment string
	Tp      model.IndexType
}

type IndexColumn struct {
	Name   string
	Length int
}

type ForeignKey struct {
	RefTable   string
	RefColumns []string
	Columns    []string
	OnDelete   int
	OnUpdate   int
}
