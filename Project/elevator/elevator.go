package elevator

import (
	elevdriver "Driver-go/elevdriver"
	"fmt"
)

// Constants matching the C defines
const (
	FloorCount           = 4
	CallButtonTypesCount = 3 // BT_HallUp, BT_HallDown, BT_Cab
)

// ElevatorBehaviour corresponds to the C enum
type ElevatorBehaviour int

const (
	EB_Idle ElevatorBehaviour = iota
	EB_DoorOpen
	EB_Moving
)

// Config holds configuration parameters
type Config struct {
	DoorOpenDurationS   float64
	ClearRequestVariant string // "All" or "InDirection" (We implement InDirection per instructions)
}

// Elevator struct (equivalent to the C struct)
type Elevator struct {
	Floor     int
	Direction elevdriver.MotorDirection
	Requests  [FloorCount][CallButtonTypesCount]bool
	Behaviour ElevatorBehaviour
	Config    Config
}

// BehaviourToString helper for logging
func (e ElevatorBehaviour) String() string {
	switch e {
	case EB_Idle:
		return "EB_Idle"
	case EB_DoorOpen:
		return "EB_DoorOpen"
	case EB_Moving:
		return "EB_Moving"
	default:
		return "EB_UNDEFINED"
	}
}

// DirectionToString helper
func DirectionToString(d elevdriver.MotorDirection) string {
	switch d {
	case elevdriver.MD_Up:
		return "MD_Up"
	case elevdriver.MD_Down:
		return "MD_Down"
	case elevdriver.MD_Stop:
		return "MD_Stop"
	default:
		return "MD_UNDEFINED"
	}
}

// Print prints the elevator state (ASCII art style like the C code)
func (e Elevator) Print() {
	fmt.Println("  +--------------------+")
	fmt.Printf("  |floor = %-2d          |\n", e.Floor)
	fmt.Printf("  |Direction  = %-12.12s|\n", DirectionToString(e.Direction))
	fmt.Printf("  |behav = %-12.12s|\n", e.Behaviour.String())
	fmt.Println("  +--------------------+")
	fmt.Println("  |  | up  | dn  | cab |")
	for f := FloorCount - 1; f >= 0; f-- {
		fmt.Printf("  | %d", f)
		for btn := 0; btn < CallButtonTypesCount; btn++ {
			if (f == FloorCount-1 && btn == int(elevdriver.BT_HallUp)) ||
				(f == 0 && btn == int(elevdriver.BT_HallDown)) {
				fmt.Print("|     ")
			} else {
				if e.Requests[f][btn] {
					fmt.Print("|  #  ")
				} else {
					fmt.Print("|  -  ")
				}
			}
		}
		fmt.Println("|")
	}
	fmt.Println("  +--------------------+")
}

// -- Requests Logic (from requests.c) --

func (e Elevator) hasRequestsAbove() bool {
	for f := e.Floor + 1; f < FloorCount; f++ {
		for btn := 0; btn < CallButtonTypesCount; btn++ {
			if e.Requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func (e Elevator) hasRequestsBelow() bool {
	for f := 0; f < e.Floor; f++ {
		for btn := 0; btn < CallButtonTypesCount; btn++ {
			if e.Requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func (e Elevator) hasRequestsHere() bool {
	for btn := 0; btn < CallButtonTypesCount; btn++ {
		if e.Requests[e.Floor][btn] {
			return true
		}
	}
	return false
}

// DirectionBehaviourPair is a return type for ChooseDirection
type DirectionBehaviourPair struct {
	Direction elevdriver.MotorDirection
	Behaviour ElevatorBehaviour
}

// ChooseDirection implements requests_chooseDirection
func (e Elevator) ChooseDirection() DirectionBehaviourPair {
	switch e.Direction {
	case elevdriver.MD_Up:
		if e.hasRequestsAbove() {
			return DirectionBehaviourPair{elevdriver.MD_Up, EB_Moving}
		} else if e.hasRequestsHere() {
			return DirectionBehaviourPair{elevdriver.MD_Down, EB_DoorOpen} // Intention of going down because previous if statement verified that there are no requests above
		} else if e.hasRequestsBelow() {
			return DirectionBehaviourPair{elevdriver.MD_Down, EB_Moving}
		} else {
			return DirectionBehaviourPair{elevdriver.MD_Stop, EB_Idle}
		}
	case elevdriver.MD_Down:
		if e.hasRequestsBelow() {
			return DirectionBehaviourPair{elevdriver.MD_Down, EB_Moving}
		} else if e.hasRequestsHere() {
			return DirectionBehaviourPair{elevdriver.MD_Up, EB_DoorOpen} // Intention of going up because previous if statement verified that there are no requests below
		} else if e.hasRequestsAbove() {
			return DirectionBehaviourPair{elevdriver.MD_Up, EB_Moving}
		} else {
			return DirectionBehaviourPair{elevdriver.MD_Stop, EB_Idle}
		}
	case elevdriver.MD_Stop:
		if e.hasRequestsHere() {
			return DirectionBehaviourPair{elevdriver.MD_Stop, EB_DoorOpen}
		} else if e.hasRequestsAbove() {
			return DirectionBehaviourPair{elevdriver.MD_Up, EB_Moving}
		} else if e.hasRequestsBelow() {
			return DirectionBehaviourPair{elevdriver.MD_Down, EB_Moving}
		} else {
			return DirectionBehaviourPair{elevdriver.MD_Stop, EB_Idle}
		}
	default:
		return DirectionBehaviourPair{elevdriver.MD_Stop, EB_Idle}
	}
}

// ShouldStop implements requests_shouldStop
func (e Elevator) ShouldStop() bool {
	switch e.Direction {
	case elevdriver.MD_Down:
		return e.Requests[e.Floor][elevdriver.BT_HallDown] ||
			e.Requests[e.Floor][elevdriver.BT_Cab] ||
			!e.hasRequestsBelow()
	case elevdriver.MD_Up:
		return e.Requests[e.Floor][elevdriver.BT_HallUp] ||
			e.Requests[e.Floor][elevdriver.BT_Cab] ||
			!e.hasRequestsAbove()
	case elevdriver.MD_Stop:
		fallthrough
	default:
		return true
	}
}

// ShouldClearImmediately implements requests_shouldClearImmediately
func (e Elevator) ShouldClearImmediately(btnFloor int, btnType elevdriver.ButtonType) bool {
	return e.Floor == btnFloor &&
		((e.Direction == elevdriver.MD_Up && btnType == elevdriver.BT_HallUp) ||
			(e.Direction == elevdriver.MD_Down && btnType == elevdriver.BT_HallDown) ||
			e.Direction == elevdriver.MD_Stop ||
			btnType == elevdriver.BT_Cab)
}

// ClearAtCurrentFloor implements requests_clearAtCurrentFloor
// This modifies the elevator state in place.
func (e *Elevator) ClearAtCurrentFloor() {
	e.Requests[e.Floor][elevdriver.BT_Cab] = false

	switch e.Direction {
	case elevdriver.MD_Up:
		if !e.hasRequestsAbove() && !e.Requests[e.Floor][elevdriver.BT_HallUp] {
			e.Requests[e.Floor][elevdriver.BT_HallDown] = false // No one wants to go up, therefore if call down, it can be cleared because it will immediatly be server
		}
		e.Requests[e.Floor][elevdriver.BT_HallUp] = false
	case elevdriver.MD_Down:
		if !e.hasRequestsBelow() && !e.Requests[e.Floor][elevdriver.BT_HallDown] {
			e.Requests[e.Floor][elevdriver.BT_HallUp] = false
		}
		e.Requests[e.Floor][elevdriver.BT_HallDown] = false
	case elevdriver.MD_Stop:
		fallthrough
	default:
		e.Requests[e.Floor][elevdriver.BT_HallUp] = false
		e.Requests[e.Floor][elevdriver.BT_HallDown] = false
	}
}
