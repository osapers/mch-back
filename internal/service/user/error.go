package user

var (
	WrongPasswordError = wrongPasswordError{}
)

type wrongPasswordError struct{}

func (e wrongPasswordError) Error() string {
	return "wrong password"
}

func newWrongPasswordError() error {
	return wrongPasswordError{}
}
