package connectionmanager

//go:generate go run github.com/arnavdugar/hsm/codegen -i=machine.yaml -o=machine.go
//go:generate go tool mockgen -source=machine.go -destination=machine_mock_test.go -package connectionmanager

// RetryHandler counts retries after the initial attempt. A successful connection
// resets the budget so later connection loss starts a fresh recovery cycle.
type RetryHandler struct {
	MaxRetries int
	Retries    int
}

var _ Handler = (*RetryHandler)(nil)

func (h *RetryHandler) CanRetry() (bool, error) {
	return h.Retries < h.MaxRetries, nil
}

func (h *RetryHandler) RecordRetry() error {
	h.Retries++
	return nil
}

func (h *RetryHandler) ResetRetries() error {
	h.Retries = 0
	return nil
}

// StartSession gives an explicitly started session a fresh retry budget.
func (h *RetryHandler) StartSession() error {
	return h.ResetRetries()
}
