package fizzBuzzMultithreaded

/*
You have the four functions:

printFizz that prints the word "fizz" to the console,
printBuzz that prints the word "buzz" to the console,
printFizzBuzz that prints the word "fizzbuzz" to the console, and
printNumber that prints a given integer to the console.
You are given an instance of the class FizzBuzz that has four functions: fizz, buzz, fizzbuzz and number. The same instance of FizzBuzz will be passed to four different threads:

Thread A: calls fizz() that should output the word "fizz".
Thread B: calls buzz() that should output the word "buzz".
Thread C: calls fizzbuzz() that should output the word "fizzbuzz".
Thread D: calls number() that should only output the integers.
Modify the given class to output the series [1, 2, "fizz", 4, "buzz", ...] where the ith token (1-indexed) of the series is:

"fizzbuzz" if i is divisible by 3 and 5,
"fizz" if i is divisible by 3 and not 5,
"buzz" if i is divisible by 5 and not 3, or
i if i is not divisible by 3 or 5.
Implement the FizzBuzz class:

FizzBuzz(int n) Initializes the object with the number n that represents the length of the sequence that should be printed.
void fizz(printFizz) Calls printFizz to output "fizz".
void buzz(printBuzz) Calls printBuzz to output "buzz".
void fizzbuzz(printFizzBuzz) Calls printFizzBuzz to output "fizzbuzz".
void number(printNumber) Calls printnumber to output the numbers.


*/

/*I dont know if my this approach will be correct or not but i think i can solve this using the single shared variable that will be used to send the singl to the all respective thread*/
type FizzBuzz struct {
	n            int
	fizzBuzzchan chan int
	fizzchan     chan int
	Buzzchan     chan int
	Numberchan   chan int
	shared_chan  chan int
}

func NewFizzBuzz(n int) *FizzBuzz {
	fb := &FizzBuzz{
		n:            n,
		fizzBuzzchan: make(chan int),
		fizzchan:     make(chan int),
		Buzzchan:     make(chan int),
		Numberchan:   make(chan int),
		shared_chan:  make(chan int, 1),
	}
	go fb.shared_state()
	return fb
}

func (fb *FizzBuzz) shared_state() {
	var shared_index int = 1
	//fb.shared_chan <- shared_index I have to ask question to myself why i have at one place use this and now i have removed this
	for shared_index <= fb.n {

		switch {

		case shared_index%3 == 0 && shared_index%5 == 0:

			fb.fizzBuzzchan <- shared_index

		case shared_index%3 == 0:

			fb.fizzchan <- shared_index

		case shared_index%5 == 0:

			fb.Buzzchan <- shared_index
		default:
			fb.Numberchan <- shared_index

		}
		/*things we stuck in the deadlock but after keeping this at the bottom it passed and worked properly need to know this thing better and ask myself why it worked*/
		shared_index = <-fb.shared_chan // updating the shared_index
	}
	if shared_index > fb.n {
		close(fb.fizzBuzzchan)
		close(fb.Buzzchan)
		close(fb.Numberchan)
		close(fb.fizzchan)
	}

}

// printFizz() outputs "fizz when the i is divisible by the 3 and not by 5"
func (fb *FizzBuzz) Fizz(printFizz func()) {
	for {
		val, ok := <-fb.fizzchan

		if !ok {

			return
		}

		printFizz()

		val++

		fb.shared_chan <- val
	}
}

// printBuzz() outputs "buzz"
func (fb *FizzBuzz) Buzz(printBuzz func()) {
	for {
		val, ok := <-fb.Buzzchan
		if !ok {

			return
		}
		printBuzz()
		val++
		fb.shared_chan <- val
	}
}

// printFizzBuzz() outputs "fizzbuzz"
func (fb *FizzBuzz) Fizzbuzz(printFizzBuzz func()) {
	for {
		val, ok := <-fb.fizzBuzzchan
		if !ok {

			return
		}
		printFizzBuzz()
		val++
		fb.shared_chan <- val
	}
}

// You may call global function `func printNumber(x int)`
// to output "x", where x is an integer.
func (fb *FizzBuzz) Number(printNumber func(int)) {
	for {
		val, ok := <-fb.Numberchan
		if !ok {

			return
		}
		printNumber(val)
		val++
		fb.shared_chan <- val
	}
}
