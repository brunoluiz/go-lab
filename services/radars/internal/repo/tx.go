package repo

import (
	"database/sql"
	"errors"
)

type Tx interface {
	Run(cb func(Querier) error) (errs error)
}

type TxExec struct {
	db      *sql.DB
	queries *Queries
}

func NewTxExec(db *sql.DB, q *Queries) *TxExec {
	return &TxExec{
		db:      db,
		queries: q,
	}
}

func (t *TxExec) Run(cb func(Querier) error) (errs error) {
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

var _ Tx = (*TxExec)(nil)
