package account

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/iowanobos/postgres-client/postgres"
	pq "github.com/iowanobos/postgres-client/postgres/query"
)

type Repository interface {
	List(ctx context.Context, query *pq.ListQuery) (*pq.ListResult[Account], error)
	BatchCreate(ctx context.Context, accounts []Account) error
}

type repository struct {
	query postgres.QueryManager
}

func NewRepository(client *postgres.Client) Repository {
	return &repository{
		query: client.QueryManager(),
	}
}

func (r *repository) List(ctx context.Context, query *pq.ListQuery) (*pq.ListResult[Account], error) {
	list := new(pq.ListResult[Account])

	itemsQb := r.listQueryBuilder(query.Filter, "id", "name")
	itemsQb = query.ApplyAll(itemsQb)

	if err := r.query.Select(ctx, &list.Items, itemsQb); err != nil {
		return nil, err
	}

	countQb := r.listQueryBuilder(query.Filter, "count(*)")
	var count []int64
	if err := r.query.Select(ctx, &count, countQb); err != nil {
		return nil, err
	}
	if len(count) > 0 {
		list.TotalCount = count[0]
	}

	return list, nil
}

func (r *repository) listQueryBuilder(filter squirrel.Sqlizer, columns ...string) squirrel.SelectBuilder {
	return pq.Qb().
		Select(columns...).
		From(TableNameAccount).
		Where(filter)
}

func (r *repository) BatchCreate(ctx context.Context, accounts []Account) error {
	qb := pq.Qb().
		Insert(TableNameAccount).
		Columns("name")

	for _, account := range accounts {
		qb = qb.Values(account.Name)
	}

	_, err := r.query.Insert(ctx, qb)
	return err
}
