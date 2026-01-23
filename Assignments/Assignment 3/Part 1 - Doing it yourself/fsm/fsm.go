package fsm

import (
	"Driver-go/elevio"
	"fmt"
	"project/elevator"
	"time"
)

// We need a channel to handle the door timeout, similar to the timer logic in C
var doorTimer *time.Timer

func init() {
	// Initialize the timer but stop it immediately so it doesn't fire
	doorTimer = time.NewTimer(time.Hour)
	doorTimer.Stop()
}

// SetAllLights sets the button lamps based on requests
func SetAllLights(e elevator.Elevator) {
	for f := 0; f < elevator.NumFloors; f++ {
		for btn := 0; btn < elevator.NumButtonTypes; btn++ {
			elevio.SetButtonLamp(elevio.ButtonType(btn), f, e.Requests[f][btn])
		}
	}
}

// OnInitializeBetweenFloors moves the elevator down if it starts between floors
func OnInitializeBetweenFloors(e *elevator.Elevator) {
	elevio.SetMotorDirection(elevio.MD_Down)
	e.Direction = elevio.MD_Down
	e.Behaviour = elevator.EB_Moving
}

// OnRequestButtonPress handles button presses
func OnRequestButtonPress(btnEvent elevio.ButtonEvent, e *elevator.Elevator, doorTimerReset chan<- bool) {
	fmt.Printf("\n\n%s(%d, %d)\n", "OnRequestButtonPress", btnEvent.Floor, btnEvent.ButtonCallType)
	e.Print()

	switch e.Behaviour {
	case elevator.EB_DoorOpen:
		if e.ShouldClearImmediately(btnEvent.Floor, btnEvent.ButtonCallType) {
			// Reset the timer
			doorTimerReset <- true
		} else {
			e.Requests[btnEvent.Floor][btnEvent.ButtonCallType] = true
		}

	case elevator.EB_Moving:
		e.Requests[btnEvent.Floor][btnEvent.ButtonCallType] = true

	case elevator.EB_Idle:
		e.Requests[btnEvent.Floor][btnEvent.ButtonCallType] = true
		pair := e.ChooseDirection()
		e.Direction = pair.Direction
		e.Behaviour = pair.Behaviour

		switch pair.Behaviour {
		case elevator.EB_DoorOpen:
			elevio.SetDoorOpenLamp(true)
			doorTimerReset <- true
			e.ClearAtCurrentFloor()
		case elevator.EB_Moving:
			elevio.SetMotorDirection(e.Direction)
		case elevator.EB_Idle:
			// Do nothing
		}
	}

	SetAllLights(*e)
	fmt.Println("\nNew state:")
	e.Print()
}

// OnFloorArrival handles arriving at a floor
func OnFloorArrival(newFloor int, e *elevator.Elevator, doorTimerReset chan<- bool) {
	fmt.Printf("\n\n%s(%d)\n", "OnFloorArrival", newFloor)
	e.Print()

	e.Floor = newFloor
	elevio.SetFloorIndicator(e.Floor)

	switch e.Behaviour {
	case elevator.EB_Moving:
		if e.ShouldStop() {
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevio.SetDoorOpenLamp(true)
			e.ClearAtCurrentFloor()
			doorTimerReset <- true
			SetAllLights(*e)
			e.Behaviour = elevator.EB_DoorOpen
		}
	default:
		// Should not happen if strictly following FSM, but safe to ignore
	}

	fmt.Println("\nNew state:")
	e.Print()
}

// OnDoorTimeout handles the door closing
func OnDoorTimeout(e *elevator.Elevator, doorTimerReset chan<- bool) {
	fmt.Printf("\n\n%s()\n", "OnDoorTimeout")
	e.Print()

	switch e.Behaviour {
	case elevator.EB_DoorOpen:
		pair := e.ChooseDirection()
		e.Direction = pair.Direction
		e.Behaviour = pair.Behaviour

		switch e.Behaviour {
		case elevator.EB_DoorOpen:
			doorTimerReset <- true
			e.ClearAtCurrentFloor()
			SetAllLights(*e)
		case elevator.EB_Moving, elevator.EB_Idle:
			elevio.SetDoorOpenLamp(false)
			elevio.SetMotorDirection(e.Direction)
		}

	default:
		// Should not happen
	}

	fmt.Println("\nNew state:")
	e.Print()
}
