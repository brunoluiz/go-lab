package repo

import (
	"database/sql"
	"errors"
)

type Tx struct {
	db      *sql.DB
	queries *Queries
}

func NewTx(db *sql.DB, q *Queries) *Tx {
	return &Tx{
		db:      db,
		queries: q,
	}
}

func (t *Tx) Run(cb func(Querier) error) (errs error) {
	tx, err := t.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if terr := tx.Rollback(); !errors.Is(terr, sql.ErrTxDone) {
			errs = errors.Join(errs, terr)
		}
	}()

	qtx := t.queries.WithTx(tx)
	if err = cb(qtx); err != nil {
		return err
	}

	return tx.Commit()
}
