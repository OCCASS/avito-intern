package pullrequest

import "errors"

var (
	ErrPrMerged               = errors.New("cannot reassign on merged PR")
	ErrNoCandidatesToReassign = errors.New("no active replacement candidate in team")
	ErrUserIsNotReviewer      = errors.New("reviewer is not assigned to this PR")
)
