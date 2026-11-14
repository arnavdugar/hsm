package queue

//go:generate go run github.com/arnavdugar/hsm/codegen -i=machine.yaml -o go=machine.go -o typescript=machine.ts -o mermaid=machine.md
//go:generate go tool mockgen -source=machine.go -destination=machine_mock_test.go -package queue

type QueueElement struct{}
