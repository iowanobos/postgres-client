package postgres

import "github.com/jackc/pgx/v5"

func RowScanStruct[T any](rows pgx.Rows) (T, error) {
	return pgx.CollectOneRow(rows, pgx.RowToStructByName[T])
}

func RowsScanStruct[T any](rows pgx.Rows) ([]T, error) {
	return pgx.CollectRows(rows, pgx.RowToStructByName[T])
}
