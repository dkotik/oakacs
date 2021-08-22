package oaksqlite

import (
	"context"
	"database/sql"
	"fmt"
	// "github.com/dkotik/oakacs/v1"
)

// var _ oakacs.TokenRepository = (*oakacs.TokenRepository)(nil)

type sessions struct {
	create              *sql.Stmt
	retrieve            *sql.Stmt
	updateLastRetrieved *sql.Stmt
	updateRole          *sql.Stmt
	updateValues        *sql.Stmt
	delete              *sql.Stmt
}

func (s *sessions) setup(table string, db *sql.DB) (err error) {
	// UUID           xid.ID
	// Differentiator string // to prevent session ID guessing
	// Identity       xid.ID
	// Role           xid.ID
	// Created        time.Time
	// LastRetrieved  time.Time
	// Values         map[string]interface{}

	if _, err = db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (uuid BLOB, identity BLOB, role BLOB, created INTEGER, lastretrieved INTEGER, differentiator TEXT, vals TEXT)", table)); err != nil {
		return
	}
	if s.create, err = db.Prepare(fmt.Sprintf("INSERT INTO `%s` VALUES(?,?,?,?,?,?)", table)); err != nil {
		return
	}
	if s.retrieve, err = db.Prepare(fmt.Sprintf("SELECT * FROM `%s` WHERE uuid=?", table)); err != nil {
		return
	}
	if s.updateLastRetrieved, err = db.Prepare(fmt.Sprintf("UPDATE `%s` SET lastretrieved=? WHERE uuid=?", table)); err != nil {
		return
	}
	if s.updateRole, err = db.Prepare(fmt.Sprintf("UPDATE `%s` SET role=? WHERE uuid=?", table)); err != nil {
		return
	}
	if s.updateValues, err = db.Prepare(fmt.Sprintf("UPDATE `%s` SET values=? WHERE uuid=?", table)); err != nil {
		return
	}
	if s.delete, err = db.Prepare(fmt.Sprintf("DELETE FROM `%s` WHERE uuid=?", table)); err != nil {
		return
	}
	return nil
}

func (s *sessions) Create(ctx context.Context) error {
	// UUID           xid.ID
	// Differentiator string // to prevent session ID guessing
	// Identity       xid.ID
	// Role           xid.ID
	// Created        time.Time
	// LastRetrieved  time.Time
	// Values         map[string]interface{}

	return nil
}