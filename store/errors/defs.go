package errors

type Error string

const (
	ErrNotExist Error = "the certificate does not exist"

)

func (err Error) Error() string {
	return string(err)
}

