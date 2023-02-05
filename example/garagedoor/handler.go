package garagedoor

import "fmt"

//go:generate go run github.com/arnavdugar/hsm/codegen -i=machine.yaml -o=machine.go
//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source=machine.go -destination=machine_mock.go -package garagedoor

type MotorDirection string

const (
	MotorDirectionStopped MotorDirection = "stop"
	MotorDirectionUp      MotorDirection = "up"
	MotorDirectionDown    MotorDirection = "down"
)

type State struct {
	MotorDirection MotorDirection
}

type GarageDoorHandler struct{}

var _ Handler = &GarageDoorHandler{}

func (handler *GarageDoorHandler) HandleButtonWhenStoppedClosing() error {
	fmt.Printf("setting motor direction: %s", MotorDirectionUp)
	return nil
}

func (handler *GarageDoorHandler) HandleButtonWhenOpening() error {
	fmt.Printf("setting motor direction: %s", MotorDirectionStopped)
	return nil
}

func (handler *GarageDoorHandler) HandleOpened() error {
	fmt.Printf("setting motor direction: %s", MotorDirectionStopped)
	return nil
}

func (handler *GarageDoorHandler) HandleButtonWhenStoppedOpening() error {
	fmt.Printf("setting motor direction: %s", MotorDirectionDown)
	return nil
}

func (handler *GarageDoorHandler) HandleButtonWhenClosing() error {
	fmt.Printf("setting motor direction: %s", MotorDirectionStopped)
	return nil
}

func (handler *GarageDoorHandler) HandleSensor() error {
	fmt.Printf("setting motor direction: %s", MotorDirectionUp)
	return nil
}

func (handler *GarageDoorHandler) HandleClosed() error {
	fmt.Printf("setting motor direction: %s", MotorDirectionStopped)
	return nil
}
