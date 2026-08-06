package lifecycle


type Transition struct {


	From string


	To string


	EventType string

}

var DefaultTransitions = []Transition{

{
	From:"CREATED",
	To:"ASSEMBLING",
	EventType:"START_ASSEMBLY",
},


{
	From:"ASSEMBLING",
	To:"TESTING",
	EventType:"START_TEST",
},


{
	From:"TESTING",
	To:"PACKAGED",
	EventType:"TEST_PASS",
},


{
	From:"PACKAGED",
	To:"SHIPPED",
	EventType:"SHIP",
},


{
	From:"SHIPPED",
	To:"DELIVERED",
	EventType:"DELIVER",
},


}