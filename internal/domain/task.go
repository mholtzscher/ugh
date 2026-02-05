package domain

import "fmt"

const (
	TaskStateInbox   = "inbox"
	TaskStateNow     = "now"
	TaskStateWaiting = "waiting"
	TaskStateLater   = "later"
	TaskStateDone    = "done"

	TaskStatesUsage = TaskStateInbox + "|" + TaskStateNow + "|" + TaskStateWaiting + "|" + TaskStateLater + "|" + TaskStateDone
)

const (
	DateLayoutYYYYMMDD = "2006-01-02"
	DateTextYYYYMMDD   = "YYYY-MM-DD"

	MetaSeparatorColon = ":"
	MetaTextKeyValue   = "key:value"
)

func IsTaskState(value string) bool {
	switch value {
	case TaskStateInbox, TaskStateNow, TaskStateWaiting, TaskStateLater, TaskStateDone:
		return true
	default:
		return false
	}
}

func InvalidStateExpectedError(value string) error {
	return fmt.Errorf("invalid state %q (expected %s)", value, TaskStatesUsage)
}

func InvalidStateMustBeError(value string) error {
	return fmt.Errorf("invalid state %q: must be %s", value, TaskStatesUsage)
}

func InvalidDateFormatError(value string) error {
	return fmt.Errorf("invalid date format: %s (expected %s)", value, DateTextYYYYMMDD)
}

func InvalidDueOnFormatError(value string) error {
	return fmt.Errorf("invalid due_on %q: expected %s", value, DateTextYYYYMMDD)
}

func InvalidMetaFormatError(value string) error {
	return fmt.Errorf("invalid meta format: %s (expected %s)", value, MetaTextKeyValue)
}
