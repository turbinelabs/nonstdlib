package executor

type Try interface {
	IsReturn() bool
	IsError() bool

	Get() interface{}
	Error() error
}

func NewTry(i interface{}, err error) Try {
	if err != nil {
		return &errorT{err}
	} else {
		return &returnT{i}
	}
}

func NewReturn(i interface{}) Try {
	return &returnT{i}
}

func NewError(err error) Try {
	return &errorT{err}
}

type Return interface {
	Try
}

type Error interface {
	Try
}

type returnT struct {
	r interface{}
}

func (r *returnT) IsReturn() bool {
	return true
}

func (r *returnT) IsError() bool {
	return false
}

func (r *returnT) Get() interface{} {
	return r.r
}

func (r *returnT) Error() error {
	panic("this Try is not an Error")
}

type errorT struct {
	e error
}

func (e *errorT) IsReturn() bool {
	return false
}

func (e *errorT) IsError() bool {
	return true
}

func (e *errorT) Get() interface{} {
	panic("this Try is not a Return")
}

func (e *errorT) Error() error {
	return e.e
}
