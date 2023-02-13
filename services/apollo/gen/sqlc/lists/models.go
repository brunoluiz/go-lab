// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package lists

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type TaskStatus string

const (
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusPending   TaskStatus = "pending"
)

func (e *TaskStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = TaskStatus(s)
	case string:
		*e = TaskStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for TaskStatus: %T", src)
	}
	return nil
}

type NullTaskStatus struct {
	TaskStatus TaskStatus
	Valid      bool // Valid is true if TaskStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullTaskStatus) Scan(value interface{}) error {
	if value == nil {
		ns.TaskStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.TaskStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullTaskStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.TaskStatus, nil
}

type List struct {
	ID        int32     `json:"-"`
	UID       string    `json:"uid"`
	Title     string    `json:"title"`
	Position  int32     `json:"position"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Task struct {
	ID        int32      `json:"id"`
	UID       string     `json:"uid"`
	Title     string     `json:"title"`
	Position  int32      `json:"position"`
	Status    TaskStatus `json:"status"`
	UpdatedAt time.Time  `json:"updated_at"`
	CreatedAt time.Time  `json:"created_at"`
}
