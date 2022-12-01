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
	itemsQb := r.listQueryBuilder(query.Filter, "id", "name")
	itemsQb = query.ApplyIteration(itemsQb)
	itemsQb = query.ApplySort(itemsQb)
	res, err := r.query.Select(ctx, itemsQb)
	if err != nil {
		return nil, err
	}

	items, err := postgres.RowsScanStruct[Account](res)
	if err != nil {
		return nil, err
	}

	countQb := r.listQueryBuilder(query.Filter, "count(*)")
	res, err = r.query.Select(ctx, countQb)
	if err != nil {
		return nil, err
	}

	var totalCount int64
	for res.Next() {
		if err = res.Scan(&totalCount); err != nil {
			return nil, err
		}
	}

	return &pq.ListResult[Account]{
		Items:      items,
		TotalCount: totalCount,
	}, nil
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
