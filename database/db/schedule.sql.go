// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: schedule.sql

package db

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

const createSchedule = `-- name: CreateSchedule :exec

INSERT INTO schedule (doctor_uuid, date) VALUES (?, ?)
`

type CreateScheduleParams struct {
	DoctorUuid uuid.UUID
	Date       time.Time
}

func (q *Queries) CreateSchedule(ctx context.Context, arg CreateScheduleParams) error {
	_, err := q.db.ExecContext(ctx, createSchedule, arg.DoctorUuid, arg.Date)
	return err
}

const deleteSchedule = `-- name: DeleteSchedule :exec

DELETE FROM schedule WHERE id = ?
`

func (q *Queries) DeleteSchedule(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteSchedule, id)
	return err
}

const listDoctorSchedule = `-- name: ListDoctorSchedule :many

SELECT id, doctor_uuid, date, created_at, updated_at FROM schedule WHERE doctor_uuid = ?
`

func (q *Queries) ListDoctorSchedule(ctx context.Context, doctorUuid uuid.UUID) ([]Schedule, error) {
	rows, err := q.db.QueryContext(ctx, listDoctorSchedule, doctorUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Schedule
	for rows.Next() {
		var i Schedule
		if err := rows.Scan(
			&i.ID,
			&i.DoctorUuid,
			&i.Date,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
