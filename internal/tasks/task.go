package tasks

import (
	"github.com/google/uuid"
	"time"
)

type Task struct {
	Id           uuid.UUID `db:"id"`
	CreationTime time.Time `db:"creation_time"`
	UserId       uuid.UUID `db:"user_id"`
	Result       []byte    `db:"result"`
}
