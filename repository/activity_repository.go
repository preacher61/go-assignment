package repository

import (
	"context"
	"database/sql"

	"preacher61/go-assignment/model"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// ActivityRepository exposes methods for performing
// operations on the `activity_logs` table.
type ActivityRepository struct {
	db *sql.DB
}

// NewActivityRepository returns a new ActivityRepository.
func NewActivityRepository() *ActivityRepository {
	db, err := OpenPgSQL()
	if err != nil {
		log.Fatal().Err(err).Msg("init pgsql failed")
	}
	return &ActivityRepository{
		db: db,
	}
}

// InsertActivities inserts multiple activities into table `activity_logs`.
func (a *ActivityRepository) InsertActivities(ctx context.Context, data []*model.Activity) error {
	txn, err := a.db.Begin()
	if err != nil {
		return errors.Wrap(err, "txn begin")
	}

	stmt, err := txn.Prepare(pq.CopyIn(tableActivityLogs, "key", "activity"))
	if err != nil {
		return errors.Wrap(err, "prepare")
	}

	for _, val := range data {
		_, err = stmt.Exec(val.Key, val.Activity)
		if err != nil {
			return errors.Wrap(err, "exec")
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return errors.Wrap(err, "exec statement")
	}

	err = stmt.Close()
	if err != nil {
		return errors.Wrap(err, "close statement")
	}

	err = txn.Commit()
	if err != nil {
		return errors.Wrap(err, "txn commit")
	}
	return nil
}
