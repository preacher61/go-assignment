package repository

import (
	"context"
	"database/sql"

	"preacher61/go-assignment/model"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const uniqActivitesCountQuery = `SELECT key, COUNT(*) AS COUNT
									FROM activity_logs
										GROUP BY key;`

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

type UniqCountResponse struct {
	Key   string `json:"key"`
	Count int64  `json:"count"`
}

func (a *ActivityRepository) GetUniqActivitesCount(ctx context.Context) {
	var res []UniqCountResponse

	rows, err := a.db.Query(uniqActivitesCountQuery)
	if err != nil {
		log.Error().Err(err).Msg("error occcured while fetching uniq activities")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var row UniqCountResponse
		if err := rows.Scan(&row.Key, &row.Count); err != nil {
			log.Error().Err(err).Msg("error occcured while fetching uniq activities")
			return
		}
		res = append(res, row)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("error occcured while fetching uniq activities")
		return
	}
	log.Info().Interface("result", res).Msg("Uniq Activities count fetched")
	return
}
