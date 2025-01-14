package mtapp

import "github.com/google/uuid"

type ProcessID uuid.UUID

func (id ProcessID) String() string {
	return id.String()
}

type ThreadID string

func (id ThreadID) String() string {
	return id.String()
}
