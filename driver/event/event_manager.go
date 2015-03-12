package event_manager

type event
const (
	floorReachedEvent
	buttonPressedEvent
	
)

func stateMachine(Event event){
	switch(Event){
		case floorReachedEvent{
			handleFloorReachedEvent()
		}
		case buttonPressedEvent{
			handleButtonPressedEvent()
		}
	}
}

func handleFloorReachedEvent(){


}

func handleButtonPressedEvent(){

}

