package oaksqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/dkotik/oakacs/v1"
	"github.com/rs/xid"
)

var _ oakacs.IntegrityLockRepository = (*locks)(nil)

type locks struct {
	db     *sql.DB
	create *sql.Stmt
	delete *sql.Stmt
	clean  *sql.Stmt
}

func (l *locks) setup(table string, db *sql.DB) (err error) {
	// Lock(context.Context, ...xid.ID) error // requires unique constraint on the table
	// Unlock(context.Context, ...xid.ID) error
	l.db = db
	if _, err = db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (id BLOB UNIQUE, deadline INTEGER)", table)); err != nil {
		return
	}
	if l.create, err = db.Prepare(fmt.Sprintf("INSERT INTO `%s` VALUES(?,?)", table)); err != nil {
		return
	}
	if l.delete, err = db.Prepare(fmt.Sprintf("DELETE FROM `%s` WHERE id=?", table)); err != nil {
		return
	}
	if l.clean, err = db.Prepare(fmt.Sprintf("DELETE FROM `%s` WHERE deadline<?", table)); err != nil {
		return
	}
	return nil
}

func (l *locks) Lock(ctx context.Context, ids ...xid.ID) error {
	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, id := range ids {
		create := tx.StmtContext(ctx, l.create)
		if _, err := create.ExecContext(ctx, id, time.Now().Add(time.Hour)); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (l *locks) Unlock(ctx context.Context, ids ...xid.ID) error {
	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, id := range ids {
		delete := tx.StmtContext(ctx, l.delete)
		if _, err := delete.ExecContext(ctx, id); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (l *locks) Clean(ctx context.Context, deadline time.Time) (int64, error) {
	result, err := l.clean.ExecContext(ctx, deadline)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
