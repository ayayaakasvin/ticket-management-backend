package helper

import "github.com/google/uuid"

func MustUUIDV7() uuid.UUID {
	v7, _ := uuid.NewV7()
	return v7
}