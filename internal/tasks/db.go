package tasks

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TasksDB interface {
	Create(*Task) error
	GetById(id uuid.UUID) (*Task, error)
	Update(*Task) error
}

type SQLTasksDB struct {
	db *sqlx.DB
}

func NewTasksDB(db *sqlx.DB) TasksDB {
	return &SQLTasksDB{db}
}

func (t *SQLTasksDB) Create(task *Task) error {
	sqlQuery := `INSERT INTO "tasks" (id, creation_time, user_id, result) 
				VALUES (:id, :creation_time, :user_id, :result)`
	_, err := t.db.NamedExec(sqlQuery, task)
	return err
}

func (t *SQLTasksDB) GetById(id uuid.UUID) (*Task, error) {
	task := &Task{}
	sqlQuery := `SELECT * FROM "tasks" WHERE id = $1`
	err := t.db.Get(task, sqlQuery, id)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (t *SQLTasksDB) Update(task *Task) error {
	_, err := t.db.NamedExec(`UPDATE tasks
									SET result = :result,
									WHERE id = :id`, task)
	return err
}
