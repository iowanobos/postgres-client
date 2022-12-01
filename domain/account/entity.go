package account

const TableNameAccount = "accounts"

type Account struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}
