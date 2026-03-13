package data

import (
	"context"
	"database/sql"
)

type Account struct {
	ID             int64
	Name           string
	Description    string
	CreatedAt      string
	InitialBalance float64
	Type           string
}

type AccountsDB struct {
	db  *sql.DB
	ctx context.Context
}

var AccountsStore *AccountsDB

func NewAccountsDB(db *sql.DB, ctx context.Context) *AccountsDB {
	return &AccountsDB{db: db, ctx: ctx}
}

func InitAccountsDB(db *sql.DB, ctx context.Context) {
	AccountsStore = NewAccountsDB(db, ctx)
}

func (a *AccountsDB) List() ([]Account, error) {
	rows, rowsErr := a.db.QueryContext(
		a.ctx,
		`SELECT id, name, description, strftime('%Y-%m-%d %H:%M:%S', created_at), initial_balance, type FROM accounts ORDER BY created_at DESC, id DESC`,
	)
	if rowsErr != nil {
		return nil, rowsErr
	}
	defer rows.Close()

	accounts := make([]Account, 0)
	for rows.Next() {
		account := Account{}
		if scanErr := rows.Scan(
			&account.ID,
			&account.Name,
			&account.Description,
			&account.CreatedAt,
			&account.InitialBalance,
			&account.Type,
		); scanErr != nil {
			return nil, scanErr
		}

		accounts = append(accounts, account)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, rowsErr
	}

	return accounts, nil
}

func (a *AccountsDB) Create(name string, description string, initialBalance float64, accountType string) error {
	_, err := a.db.ExecContext(
		a.ctx,
		`INSERT INTO accounts(name, description, initial_balance, type) VALUES(?, ?, ?, ?)`,
		name,
		description,
		initialBalance,
		accountType,
	)

	return err
}

func (a *AccountsDB) Update(id int64, name string, description string, initialBalance float64, accountType string) error {
	_, err := a.db.ExecContext(
		a.ctx,
		`UPDATE accounts SET name = ?, description = ?, initial_balance = ?, type = ? WHERE id = ?`,
		name,
		description,
		initialBalance,
		accountType,
		id,
	)

	return err
}

func (a *AccountsDB) Delete(id int64) error {
	_, err := a.db.ExecContext(a.ctx, `DELETE FROM accounts WHERE id = ?`, id)
	return err
}
