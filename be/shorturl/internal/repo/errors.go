package repo

type (
	NotFoundError     struct{}
	RecordExistsError struct{}
)

func (_ *NotFoundError) Error() string {
	return "not found"
}

func (_ *RecordExistsError) Error() string {
	return "record exists"
}

var ErrRecordExists = &RecordExistsError{}
