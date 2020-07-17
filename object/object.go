package object

import (
	"fmt"
)

type ObjectType string

const (
	INTEGER_OBJ     = "INTEGER"
	PRINT_VALUE_OBJ = "PRINT_VALUE"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

// Integer
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

// PrintValue : Print value object
type PrintValue struct {
	Value Object
}

func (pv *PrintValue) Type() ObjectType { return PRINT_VALUE_OBJ }
func (pv *PrintValue) Inspect() string  { return pv.Value.Inspect() }
