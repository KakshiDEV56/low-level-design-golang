package designundergroundsystem

/*
LeetCode 1396. Design Underground System

An underground railway system is keeping track of customer travel times between
different stations. They are using this data to calculate the average time it takes to
travel from one station to another.

Implement the UndergroundSystem class:
  - NewUndergroundSystem() initializes the object of the system.
  - CheckIn(id int, stationName string, t int) a customer with card id equal to id checks
    in at the station stationName at time t. A customer can only be checked into one
    place at a time.
  - CheckOut(id int, stationName string, t int) a customer with card id equal to id checks
    out from the station stationName at time t.
  - GetAverageTime(startStation, endStation string) float64 returns the average time it
    takes to travel from startStation to endStation. The average time is computed from
    all the previous traveling times from startStation to endStation that happened
    directly, meaning a check in at startStation followed by a check out from
    endStation. The time it takes to travel from startStation to endStation may be
    different from the time it takes to travel from endStation to startStation.
    There will be at least one journey from startStation to endStation before
    GetAverageTime is called.

Example 1:

	Input:
	["UndergroundSystem", "checkIn", "checkIn", "checkIn", "checkOut", "checkOut",
	 "checkOut", "getAverageTime", "getAverageTime", "checkIn", "getAverageTime",
	 "checkOut", "getAverageTime"]
	[[], [45, "Leyton", 3], [32, "Paradise", 8], [27, "Leyton", 10],
	 [45, "Waterloo", 15], [27, "Waterloo", 20], [32, "Cambridge", 22],
	 ["Paradise", "Cambridge"], ["Leyton", "Waterloo"], [10, "Leyton", 24],
	 ["Leyton", "Waterloo"], [10, "Waterloo", 38], ["Leyton", "Waterloo"]]
	Output:
	[null, null, null, null, null, null, null, 14.00000, 11.00000, null, 11.00000, null,
	 12.00000]

	Explanation:
	UndergroundSystem us = new UndergroundSystem();
	us.checkIn(45, "Leyton", 3);
	us.checkIn(32, "Paradise", 8);
	us.checkIn(27, "Leyton", 10);
	us.checkOut(45, "Waterloo", 15);   // customer 45 "Leyton" -> "Waterloo" in 15-3 = 12
	us.checkOut(27, "Waterloo", 20);   // customer 27 "Leyton" -> "Waterloo" in 20-10 = 10
	us.checkOut(32, "Cambridge", 22);  // customer 32 "Paradise" -> "Cambridge" in 22-8 = 14
	us.getAverageTime("Paradise", "Cambridge"); // return 14.0
	us.getAverageTime("Leyton", "Waterloo");    // return 11.0, avg of 12 and 10
	us.checkIn(10, "Leyton", 24);
	us.getAverageTime("Leyton", "Waterloo");    // return 11.0
	us.checkOut(10, "Waterloo", 38);            // customer 10 "Leyton" -> "Waterloo" in 38-24 = 14
	us.getAverageTime("Leyton", "Waterloo");    // return 12.0, avg of 12, 10 and 14

Constraints:
  - 1 <= id, t <= 10^6
  - 1 <= stationName.length, startStation.length, endStation.length <= 10
  - All station names consist of only lowercase and uppercase English letters and digits.
  - There will be at most 2*10^4 calls in total to CheckIn, CheckOut, and GetAverageTime.
  - Answers within 10^-5 of the actual value will be accepted.

TODO:
 1. Add the fields you need on UndergroundSystem below. Hint:
      - a map from customer id -> the (stationName, checkInTime) of their currently open
        trip, populated by CheckIn and consumed by CheckOut, and
      - a map keyed by the (startStation, endStation) pair -> a running total travel time
        and a trip count, so GetAverageTime is O(1) instead of re-scanning history. (Go
        maps can be keyed by a struct with comparable fields, which is a clean way to key
        on a station pair.)
 2. Implement NewUndergroundSystem: initialize your maps.
 3. Implement CheckIn: record the open trip for id.
 4. Implement CheckOut: look up id's open trip, compute the elapsed time, fold it into
    the (start, end) pair's running total/count, and remove id's open trip.
 5. Implement GetAverageTime: return runningTotal / count as a float64 for the given
    (start, end) pair.
*/

type UndergroundSystem struct {
	// TODO: add your fields here
}

func NewUndergroundSystem() *UndergroundSystem {
	// TODO: implement
	return &UndergroundSystem{}
}

func (u *UndergroundSystem) CheckIn(id int, stationName string, t int) {
	// TODO: implement
}

func (u *UndergroundSystem) CheckOut(id int, stationName string, t int) {
	// TODO: implement
}

func (u *UndergroundSystem) GetAverageTime(startStation string, endStation string) float64 {
	// TODO: implement
	return 0
}
