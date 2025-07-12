package bqb

import "strings"

// Dialect holds the Query dialect
type Dialect string

const (
	// PGSQL postgres dialect
	PGSQL Dialect = "postgres"
	// MYSQL MySQL dialect
	MYSQL Dialect = "mysql"
	// RAW dialect uses no parameter conversion
	RAW Dialect = "raw"
	// SQL generic dialect
	SQL Dialect = "sql"

	paramPh  = "{{xX_PARAM_Xx}}"
	ignorePH = "{{OOO_INVALID_OOO}}"
)

// Embedded is a string type that is directly embedded into the query.
// Note: Like Embedder, this is not to be used for untrusted input.
type Embedded string

// Embedder embeds a value directly into a query string.
// Note: Since this is embedded and not bound,
// attention must be paid to sanitizing this input.
type Embedder interface {
	RawValue() string
}

type Expanded[T any] []T

func (e Expanded[T]) ExpandItems() []any {
	return sliceToAny(e)
}

func Expand[T any](v []T) Expanded[T] {
	return v
}

type ExpandedJoin struct {
	Expander
	join func([]string) string
}

var _ Expander = (*ExpandedJoin)(nil)

func (e ExpandedJoin) ExpandJoin(values []string) string {
	return e.join(values)
}

func ExpandedJoinFunc(e Expander, join func([]string) string) ExpandedJoin {
	return ExpandedJoin{
		Expander: e,
		join:     join,
	}
}

func ExpandJoinFunc[T any](v []T, join func([]string) string) ExpandedJoin {
	return ExpandedJoinFunc(Expanded[T](v), join)
}

func ExpandJoin[T any](v []T, join string) ExpandedJoin {
	return ExpandJoinFunc(v, func(values []string) string {
		return strings.Join(values, join)
	})
}

type Expander interface {
	ExpandItems() []any
}

type ExpanderJoin interface {
	ExpandJoin(values []string) string
}

// JsonMap is a custom type which tells bqb to convert the parameter to
// a JSON object without requiring reflection.
type JsonMap map[string]any

// JsonList is a type that tells bqb to convert the parameter to a JSON
// list without requiring reflection.
type JsonList []any

// Folded is a type that tells bqb to NOT spread the list into individual
// parameters.
type Folded []any

// ToFolded converts a slice to a ValueArray.
func ToFolded[T any](slice []T) Folded {
	valueArr := make(Folded, len(slice))
	for i, v := range slice {
		valueArr[i] = v
	}
	return valueArr
}
