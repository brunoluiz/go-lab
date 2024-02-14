package repo

import (
	"database/sql"
	"errors"
)

var _ Tx = (*TxExec)(nil)

type (
	// Tx allows the user to run a callback function within a transaction
	Tx interface {
		Run(cb func(Querier) error) (errs error)
	}

	// QuerierTx *Queries object exposes .WithTx, but not in the Querier interface
	// This interface is used to catch the WithTx method
	QuerierTx interface {
		WithTx(tx *sql.Tx) *Queries
	}

	// TxExec implementation of the Tx interface
	TxExec struct {
		db      *sql.DB
		queries QuerierTx
	}
)

func NewTxExec(db *sql.DB, q QuerierTx) *TxExec {
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
