package query

import "github.com/Masterminds/squirrel"

type ListQuery struct {
	Filter    squirrel.Sqlizer
	iteration iteration
	sorts     []sort
}

func NewListQuery(filter squirrel.Sqlizer) *ListQuery {
	return &ListQuery{Filter: filter}
}

func (q *ListQuery) WithIteration(value iteration) *ListQuery {
	q.iteration = value
	return q
}

func (q *ListQuery) AddSort(field string, isDesc bool) *ListQuery {
	q.sorts = append(q.sorts, sort{field, isDesc})
	return q
}

func (q *ListQuery) ApplyAll(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
	builder = q.ApplyIteration(builder)
	return q.ApplySort(builder)
}

func (q *ListQuery) ApplyIteration(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
	if q == nil || q.iteration == nil {
		return builder
	}
	return q.iteration.applyIteration(builder)
}

func (q *ListQuery) ApplySort(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
	if q == nil || len(q.sorts) == 0 {
		return builder
	}

	orderBys := make([]string, len(q.sorts))
	for i, v := range q.sorts {
		order := "asc"
		if v.IsDescOrder {
			order = "desc"
		}

		orderBys[i] = v.Field + " " + order
	}

	return builder.OrderBy(orderBys...)
}

type ListResult[T any] struct {
	Items      []T
	TotalCount int64
}

func Qb() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}
