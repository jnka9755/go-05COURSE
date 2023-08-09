package course

import (
	"errors"
	"fmt"
)

var ErrNameRequired = errors.New("name is required")
var ErrStarDateRequired = errors.New("start_date is required")
var ErrEndDateRequired = errors.New("end_date is required")

var ErrInvalidStartDate = errors.New("invalid start_date")
var ErrInvalidEndtDate = errors.New("invalid end_date")

type ErrNotFound struct {
	UserID string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("Course with ID -> '%s' doesn't exist", e.UserID)
}
