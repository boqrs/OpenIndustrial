package mes

import "errors"

var (
	ErrLineNotFound           = errors.New("production line not found")
	ErrProcessNotFound        = errors.New("process definition not found")
	ErrProcessStepNotFound    = errors.New("process step not found")
	ErrStationNotFound        = errors.New("station not found")
	ErrCapabilityNotFound     = errors.New("station capability not found")
	ErrTaskNotFound           = errors.New("station task not found")
	ErrInvalidStateTransition = errors.New("invalid task state transition")
	ErrNoAvailableStation     = errors.New("no available station with required capability")
)