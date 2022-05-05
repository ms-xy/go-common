package configuration

type WrappedError struct {
	msg   string
	inner error
}

func (e WrappedError) Error() string {
	return e.msg + ":\n" + e.inner.Error()
}

func (e WrappedError) Unwrap() error {
	return e.inner
}

func NewError(msg string, err error) error {
	return WrappedError{
		msg:   msg,
		inner: err,
	}
}
