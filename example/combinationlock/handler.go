package combinationlock

//go:generate go run github.com/arnavdugar/hsm/codegen -i=machine.yaml -o=machine.go
//go:generate go tool mockgen -source=machine.go -destination=machine_mock_test.go -package combinationlock

// LockHandler recognizes the fixed demonstration combination 4, 2, 7.
// Each guard receives the same Digit action's data at a different prefix state.
type LockHandler struct{}

var _ Handler = LockHandler{}

func (LockHandler) IsFirstDigit(digit int) (bool, error) {
	return digit == 4, nil
}

func (LockHandler) IsSecondDigit(digit int) (bool, error) {
	return digit == 2, nil
}

func (LockHandler) IsThirdDigit(digit int) (bool, error) {
	return digit == 7, nil
}
