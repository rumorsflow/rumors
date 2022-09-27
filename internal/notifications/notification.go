package notifications

type Notification interface {
	Forbidden(to any)
	Success(to any, text string)
	Error(to any, text string)
	Err(to any, err error)
	Send(to any, text string, options ...any)
}
