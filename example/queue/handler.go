package queue

//go:generate go run github.com/arnavdugar/hsm/codegen -i=machine.yaml -o=machine.go
//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source=machine.go -destination=machine_mock.go -package queue

type QueueElement struct{}
