package main

import (
	"Driver-go/elevio"
	"fmt"
	"project/elevator"
	"project/fsm"
	"time"
)

func main() {
	// 1. Initialize Hardware
	numFloors := 4
	elevio.Init("localhost:15657", numFloors)

	// 2. Initialize Elevator State
	elev := elevator.Elevator{
		Floor:     -1,
		Direction: elevio.MD_Stop,
		Behaviour: elevator.EB_Idle,
		Config: elevator.Config{
			DoorOpenDurationS: 3.0,
		},
	}

	// 3. Check if we need to move to a floor first (Startup logic)
	// Note: elevio.GetFloor() returns -1 if between floors
	if floor := elevio.GetFloor(); floor == -1 {
		fsm.OnInitializeBetweenFloors(&elev)
	} else {
		elev.Floor = floor
	}

	// 4. Create channels for IO events
	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	// 5. Timer logic
	// We use a timer that we can reset.
	// We use a buffered channel to trigger the timer reset to avoid blocking in FSM functions
	doorTimeout := time.NewTimer(time.Duration(elev.Config.DoorOpenDurationS) * time.Second)
	doorTimeout.Stop() // Don't start immediately
	doorTimerReset := make(chan bool)

	fmt.Println("Elevator started!")

	for {
		select {
		// Handle Button Press
		case btnEvent := <-drv_buttons:
			fsm.OnRequestButtonPress(btnEvent, &elev, doorTimerReset)

		// Handle Floor Arrival
		case newFloor := <-drv_floors:
			fsm.OnFloorArrival(newFloor, &elev, doorTimerReset)

		// Handle Door Timeout
		case <-doorTimeout.C:
			fsm.OnDoorTimeout(&elev, doorTimerReset)

		// Handle Timer Reset
		case <-doorTimerReset:
			// Drain the channel if it has expired but not yet received
			if !doorTimeout.Stop() {
				select {
				case <-doorTimeout.C:
				default:
				}
			}
			doorTimeout.Reset(time.Duration(elev.Config.DoorOpenDurationS) * time.Second)

		// Handle Obstruction (Optional/Extra)
		case obstr := <-drv_obstr:
			// fmt.Printf("Obstruction: %v\n", obstr)
			// // Implementation depends on specific requirements, often pauses the door timer
			// if obstr && elev.Behaviour == elevator.EB_DoorOpen {
			// 	elevio.SetDoorOpenLamp(true) // Ensure it's on
			// 	if !doorTimeout.Stop() {
			// 		select {
			// 		case <-doorTimeout.C:
			// 		default:
			// 		}
			// 	}
			// 	// Timer is effectively paused (stopped).
			// 	// Logic to resume when obstruction clears would go here.
			// } else if !obstr && elev.Behaviour == elevator.EB_DoorOpen {
			// 	// Resume timer
			// 	doorTimeout.Reset(time.Duration(elev.Config.DoorOpenDurationS) * time.Second)
			// }

		// Handle Stop Button (Optional/Extra)
		case stop := <-drv_stop:
			// fmt.Printf("Stop button: %v\n", stop)
			// for f := 0; f < numFloors; f++ {
			// 	for b := elevio.ButtonType(0); b < 3; b++ {
			// 		elevio.SetButtonLamp(b, f, false)
			// 	}
			// }
		}
	}
}
