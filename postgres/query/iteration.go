package query

import "github.com/Masterminds/squirrel"

type iteration interface {
	applyIteration(builder squirrel.SelectBuilder) squirrel.SelectBuilder
}

type OffsetIteration struct {
	Limit  uint64
	Offset uint64
}

func (v OffsetIteration) applyIteration(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
	return builder.
		Limit(v.Limit).
		Offset(v.Offset)
}

type LastIDIteration struct {
	Limit  uint64
	Field  string
	LastID any
}

func (v LastIDIteration) applyIteration(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
	return builder.
		OrderBy(v.Field).
		Limit(v.Limit).
		Where(squirrel.Gt{v.Field: v.LastID})
}

type Pagination struct {
	Number uint64
	Size   uint64
}

func (v Pagination) applyIteration(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
	return builder.
		Limit(v.Size).
		Offset(v.Size * (v.Number - 1)) // TODO: Если считать страницы с 1
}
