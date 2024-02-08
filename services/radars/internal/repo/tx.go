package repo

import (
	"database/sql"
	"errors"
)

type Tx func(func(Querier) error) error

func NewTx(db *sql.DB, q *Queries) func(func(Querier) error) error {
	return func(cb func(Querier) error) (errs error) {
		tx, err := db.Begin()
		if err != nil {
			return err
		}

		defer func() {
			if terr := tx.Rollback(); !errors.Is(terr, sql.ErrTxDone) {
				errs = errors.Join(errs, terr)
			}
		}()

		qtx := q.WithTx(tx)
		if err = cb(qtx); err != nil {
			return err
		}

		return tx.Commit()
	}
}
