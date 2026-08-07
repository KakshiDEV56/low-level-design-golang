package designparkingsystem

/*
LeetCode 1603. Design Parking System

Design a parking system for a parking lot. The parking lot has three kinds of parking
spaces: big, medium, and small, with a fixed number of slots for each size.

Implement the ParkingSystem class:
  - NewParkingSystem(big, medium, small int) initializes a ParkingSystem object with the
    given number of slots for each parking space.
  - AddCar(carType int) bool checks whether there is a parking space of carType for the
    car that wants to get into the parking lot. carType can be of three kinds: big (1),
    medium (2), and small (3). A car can only park in a parking space of its carType.
    If there is no space available, returns false, else parks the car in that size space
    and returns true.

Example 1:

	Input:
	["ParkingSystem", "addCar", "addCar", "addCar", "addCar"]
	[[1, 1, 0], [1], [2], [3], [1]]
	Output:
	[null, true, true, false, false]

	Explanation:
	ParkingSystem ps = new ParkingSystem(1, 1, 0);
	ps.addCar(1); // return true, there's 1 available slot for a big car
	ps.addCar(2); // return true, there's 1 available slot for a medium car
	ps.addCar(3); // return false, there's no available slot for a small car
	ps.addCar(1); // return false, there's no available slot for a big car; it's already
	              // occupied.

Constraints:
  - 0 <= big, medium, small <= 1000
  - carType is 1, 2, or 3.
  - At most 1000 calls will be made to AddCar.

TODO:
 1. Add the fields you need on ParkingSystem below. This one is deliberately simple —
    three integer counters (remaining slots per size) are enough.
 2. Implement NewParkingSystem: store the initial slot counts for big, medium, small.
 3. Implement AddCar: based on carType (1=big, 2=medium, 3=small), check whether that
    size still has a remaining slot. If so, decrement the counter and return true.
    Otherwise return false. This is a good spot to practice a clean switch statement
    rather than three near-identical if-blocks.
*/

type ParkingSystem struct {
	// TODO: add your fields here
}

func NewParkingSystem(big int, medium int, small int) *ParkingSystem {
	// TODO: implement
	return &ParkingSystem{}
}

func (p *ParkingSystem) AddCar(carType int) bool {
	// TODO: implement
	return false
}
