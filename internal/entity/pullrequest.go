package entity

import (
	"fmt"
	"slices"
)

type Status string

const (
	StatusOpen   Status = "OPEN"
	StatusMerged Status = "MERGED"
)

var allowedStatuses = []Status{
	StatusOpen,
	StatusMerged,
}

func NewStatus(value string) (Status, error) {
	s := Status(value)
	if !slices.Contains(allowedStatuses, s) {
		return "", fmt.Errorf("status %s is incorrect", value)
	}
	return s, nil
}

type PullRequest struct {
	Id     string
	Name   string
	Author User
	Status Status
}
