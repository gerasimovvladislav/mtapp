package mtapp

import "github.com/google/uuid"

type ProcessID uuid.UUID

func (pid ProcessID) String() string {
	return pid.String()
}

type ThreadID string

func (tid ThreadID) String() string {
	return string(tid)
}
