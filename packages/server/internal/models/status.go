package models

type Status string

const (
	StatusDone  Status = "done"
	StatusRetry Status = "retry"
	StatusDLQ   Status = "dlq"
)
