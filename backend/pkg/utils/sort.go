package utils

import (
	"time"

	"github.com/lukeelten/micc-creation-day-26/backend/pkg/models"
)

func SortTimeDesc(a time.Time, b time.Time) int {
	if a.After(b) {
		return 1
	}
	if a.Before(b) {
		return -1
	}
	return 0
}

func SortStatesByCompletedDesc(a, b *models.StatesRecord) int {
	if a.Completed == nil && b.Completed != nil {
		return 1
	}

	if b.Completed == nil && a.Completed != nil {
		return -1
	}

	if a.Completed == nil && b.Completed == nil {
		return 0
	}

	return SortTimeDesc(*b.Completed, *a.Completed)
}
