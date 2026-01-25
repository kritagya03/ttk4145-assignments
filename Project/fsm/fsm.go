package fsm

import (
	"Driver-go/elevdriver"
	"fmt"
	"project/elevator"
	"time"
)

// We need a channel to handle the door timeout, similar to the timer logic in C
var doorTimer *time.Timer

// func init() {
// 	fmt.Println("init: initialized one hour timer, but don't start it.")
// 	// Initialize the timer but stop it immediately so it doesn't fire
// 	doorTimer = time.NewTimer(time.Hour)
// 	doorTimer.Stop()
// }

// SetAllLights sets the button lamps based on requests
func SetAllLights(e elevator.Elevator) {
	for f := 0; f < elevator.FloorCount; f++ {
		for btn := 0; btn < elevator.CallButtonTypesCount; btn++ {
			elevdriver.SetButtonLamp(elevdriver.ButtonType(btn), f, e.Requests[f][btn])
		}
	}
	fmt.Println("SetAllLights: updated all lights from e.Requests.")
}

// OnInitializeBetweenFloors moves the elevator down if it starts between floors
func OnInitializeBetweenFloors(e *elevator.Elevator) {
	fmt.Println("OnInitializeBetweenFloors: wanting to move down.")
	elevdriver.SetMotorDirection(elevdriver.MD_Down)
	e.Direction = elevdriver.MD_Down
	e.Behaviour = elevator.EB_Moving
}

// OnRequestButtonPress handles button presses
func OnRequestButtonPress(btnEvent elevdriver.ButtonEvent, e *elevator.Elevator, doorTimerReset chan<- bool) {
	fmt.Printf("\n\n%s(%d, %d)\n", "OnRequestButtonPress", btnEvent.Floor, btnEvent.ButtonCallType)
	e.Print()

	switch e.Behaviour {
	case elevator.EB_DoorOpen:
		if e.ShouldClearImmediately(btnEvent.Floor, btnEvent.ButtonCallType) {
			fmt.Println("OnRequestButtonPress - EB_DoorOpen - ShouldClearImmediately=true: Wanting to start the timer again.")
			// Reset the timer
			doorTimerReset <- true
		} else {
			e.Requests[btnEvent.Floor][btnEvent.ButtonCallType] = true
			fmt.Println("OnRequestButtonPress - EB_DoorOpen - ShouldClearImmediately=false: Added order to e.Requests.")
		}

	case elevator.EB_Moving:
		e.Requests[btnEvent.Floor][btnEvent.ButtonCallType] = true
		fmt.Println("OnRequestButtonPress - EB_Moving: Added order to e.Requests.")

	case elevator.EB_Idle:
		e.Requests[btnEvent.Floor][btnEvent.ButtonCallType] = true
		fmt.Println("OnRequestButtonPress - EB_Idle: Added order to e.Requests.")
		directionBehaviourPair := e.ChooseDirection()
		e.Direction = directionBehaviourPair.Direction
		e.Behaviour = directionBehaviourPair.Behaviour
		fmt.Println("OnRequestButtonPress - EB_Idle: Chosen new direction:", directionBehaviourPair)

		switch directionBehaviourPair.Behaviour {
		case elevator.EB_DoorOpen:
			elevdriver.SetDoorOpenLamp(true)
			doorTimerReset <- true
			e.ClearAtCurrentFloor()
			fmt.Println("OnRequestButtonPress - EB_Idle - New behaviour is EB_DoorOpen: Opened door, wanting to start the timer, maybe remove order(s) at floor if possible.")
		case elevator.EB_Moving:
			elevdriver.SetMotorDirection(e.Direction)
			fmt.Println("OnRequestButtonPress - EB_Idle - New behaviour is EB_Moving: set the motor direction: ", e.Direction)
		case elevator.EB_Idle:
			// Do nothing
			fmt.Println("OnRequestButtonPress - EB_Idle - New behaviour is EB_Idle: Doing nothing")
		}
	}

	fmt.Println("OnRequestButtonPress: Wanting to update all lights.")
	SetAllLights(*e)
	fmt.Println("\nNew state:")
	e.Print()
}

// OnFloorArrival handles arriving at a floor
func OnFloorArrival(newFloor int, e *elevator.Elevator, doorTimerReset chan<- bool) {
	fmt.Printf("\n\n%s(%d)\n", "OnFloorArrival", newFloor)
	e.Print()

	fmt.Println("OnFloorArrival: wanting to set floor indicator.")
	e.Floor = newFloor
	elevdriver.SetFloorIndicator(e.Floor)

	switch e.Behaviour {
	case elevator.EB_Moving:
		if e.ShouldStop() {
			fmt.Println("OnFloorArrival - EB_Moving - e.ShouldStop()==True: wanting to stop, open door, maybe clear order(s) at floor, reset timer, update lights.")
			elevdriver.SetMotorDirection(elevdriver.MD_Stop)
			elevdriver.SetDoorOpenLamp(true)
			e.ClearAtCurrentFloor()
			doorTimerReset <- true
			SetAllLights(*e)
			e.Behaviour = elevator.EB_DoorOpen
		}
	default:
		// Can enter here if initializing elevator on a floor
		fmt.Printf("OnFloorArrival - e.Behaviour==default(%v): Doing nothing.\n", e.Behaviour)
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
		fmt.Println("OnDoorTimeout - EB_DoorOpen: choosing new direction.")
		directionBehaviourPair := e.ChooseDirection()
		e.Direction = directionBehaviourPair.Direction
		e.Behaviour = directionBehaviourPair.Behaviour

		switch e.Behaviour {
		case elevator.EB_DoorOpen:
			fmt.Println("OnDoorTimeout - EB_DoorOpen - new behaviour is EB_DoorOpen: wanting to reset timer, maybe clear order(s) at floor, updating lights.")
			doorTimerReset <- true
			e.ClearAtCurrentFloor()
			SetAllLights(*e)
		case elevator.EB_Moving, elevator.EB_Idle:
			fmt.Printf("OnDoorTimeout - EB_DoorOpen - new behaviour is %v: opening door, setting direction %v.\n", e.Behaviour, e.Direction)
			elevdriver.SetDoorOpenLamp(false)
			elevdriver.SetMotorDirection(e.Direction)
		}

	default:
		// If initializing the elevator on a floor, this happens
		fmt.Printf("OnDoorTimeout - e.Behaviour==default(%v): Doing nothing.\n", e.Behaviour)
		// Should not happen
	}

	fmt.Println("\nNew state:")
	e.Print()
}
