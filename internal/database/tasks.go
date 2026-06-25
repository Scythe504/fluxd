package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (t *Task) Scan(row pgx.Row) error {
	return row.Scan(
		&t.ID, &t.PayloadSlug, &t.Payload, &t.RetryCount,
		&t.MaxRetryCount, &t.LastError, &t.ExecutionScheduleTime,
		&t.ExecutionIntervalSeconds, &t.CronExpression, &t.TaskType,
		&t.Status, &t.AllocatedUnit, &t.CreatedAt, &t.UpdatedAt, &t.DeletedAt,
	)
}

func (s *service) GetTask(ctx context.Context) (Task, error) {
	query := `UPDATE tasks
			SET status = $1
			WHERE id = (
				SELECT id
				FROM tasks
				WHERE status = 'queued' AND deleted_at IS NULL
				ORDER BY created_at ASC
				LIMIT 1
				FOR UPDATE SKIP LOCKED
			)
			RETURNING id, payload_slug, payload, retry_count, max_retry_count, 
				last_error, execution_schedule_time, execution_interval_seconds,
				cron_expression, task_type, status, allocated_unit, 
				created_at, updated_at, deleted_at
		`
	row := s.pool.QueryRow(ctx, query, TaskStatusRunning)

	task := Task{}

	if err := task.Scan(row); err != nil {
		return task, err
	}

	return task, nil
}
