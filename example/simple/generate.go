package simple

//go:generate go run github.com/arnavdugar/hsm/codegen -i=machine.yaml -o=machine.go
//go:generate go tool mockgen -source=machine.go -destination=machine_mock_test.go -package simple
