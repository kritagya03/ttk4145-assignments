package main

import (
	"Driver-go/elevio"
	"fmt"
	"slices"
)

type ElevatorState int

const (
	Startup             ElevatorState = 0
	WaitingForCalls                   = 1
	UpMotorMove                       = 2
	UpStopAndOpenDoor                 = 3
	DownMotorMove                     = 4
	DownStopAndOpenDoor               = 5
)

const floorCount int = 4

var lastFloor int = -1
var motorDirection elevio.MotorDirection

var drv_buttons chan elevio.ButtonEvent = make(chan elevio.ButtonEvent)
var drv_floors chan int = make(chan int)
var drv_obstr chan bool = make(chan bool)
var drv_stop chan bool = make(chan bool)

var currentState ElevatorState = Startup

const buttonCount int = floorCount*2 - 2 + floorCount

var upCallList []int = make([]int, 0, buttonCount)
var downCallList []int = make([]int, 0, buttonCount)
var noDirectionCallList []int = make([]int, 0, buttonCount)

func getCurrentFloor() int {
	// Returns -1 if the elevator is not at a floor
	return elevio.GetFloor()
}

func getLargestFloorInSlice(s []int) int {
	if len(s) == 0 {
		return -1
	}
	return slices.Max(s)
}

func getSmallestFloorInSlice(s []int) int {
	if len(s) == 0 {
		return -1
	}
	return slices.Min(s)
}

func updateMotorDirection(newDirection elevio.MotorDirection) {
	if motorDirection != newDirection {
		motorDirection = newDirection
		if motorDirection != elevio.MD_Stop {
			elevio.SetDoorOpenLamp(false)
		}
		elevio.SetMotorDirection(motorDirection)
	}
}

func getDirectionToFloor(targetFloor int) elevio.MotorDirection {
	if targetFloor > lastFloor {
		return elevio.MD_Up
	} else if targetFloor < lastFloor {
		return elevio.MD_Down
	} else {
		return elevio.MD_Stop
	}
}

func setupElevator() {
	elevio.Init("localhost:15657", floorCount)
	updateMotorDirection(elevio.MD_Stop)
	elevio.SetDoorOpenLamp(false)
	if getCurrentFloor() == -1 {
		updateMotorDirection(elevio.MD_Down)
	}
}

func startSourcesOnFibers() {
	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)
}

func server() {
	for {
		select {
		case buttonEvent := <-drv_buttons:
			fmt.Printf("%+v\n", buttonEvent)
			if buttonEvent.Floor == getCurrentFloor() {
				continue // Do nothing if the elevator is already at the requested floor
			}
			elevio.SetButtonLamp(buttonEvent.Button, buttonEvent.Floor, true)
			switch currentState {
			case WaitingForCalls:
				directionToFloor := getDirectionToFloor(buttonEvent.Floor)
				switch directionToFloor {
				case elevio.MD_Up:
					upCallList = append(upCallList, buttonEvent.Floor)
					updateMotorDirection(elevio.MD_Up)
					currentState = UpMotorMove
				case elevio.MD_Down:
					downCallList = append(downCallList, buttonEvent.Floor)
					updateMotorDirection(elevio.MD_Down)
					currentState = DownMotorMove
				}
			case UpMotorMove:
				if slices.Contains(upCallList, buttonEvent.Floor) {
					continue // Do nothing, already in the list
				}
				directionToFloor := getDirectionToFloor(buttonEvent.Floor)
				if 
				// if buttonEvent.Button == elevio.BT_HallUp && directionToFloor == elevio.MD_Up &&  {
			}

		case floorNumber := <-drv_floors:
			fmt.Printf("%+v\n", floorNumber)
			switch currentState {
			case Startup:
				lastFloor = floorNumber
				updateMotorDirection(elevio.MD_Stop)
				currentState = WaitingForCalls
			}
			// if a == floorCount-1 {
			// 	d = elevio.MD_Down
			// } else if a == 0 {
			// 	d = elevio.MD_Up
			// }
			// elevio.SetMotorDirection(d)

		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			// if a {
			// 	elevio.SetMotorDirection(elevio.MD_Stop)
			// } else {
			// 	elevio.SetMotorDirection(d)
			// }

		case a := <-drv_stop:
			fmt.Printf("%+v\n", a)
			// for f := 0; f < floorCount; f++ {
			// 	for b := elevio.ButtonType(0); b < 3; b++ {
			// 		elevio.SetButtonLamp(b, f, false)
			// 	}
			// }
		}
	}
}

func main() {
	setupElevator()
	startSourcesOnFibers()
	server()
}
